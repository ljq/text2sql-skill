// Copyright 2024 Text2SQL Skill Engine
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Jaco Liu (Jianqiu Liu) <ljqlab@gmail.com>
// GitHub: https://github.com/ljq

package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"text2sql-skill/config"
	"text2sql-skill/interfaces"
	"text2sql-skill/utils"
)

type Text2SQLSkill struct {
	mu             sync.Mutex
	db             *sql.DB
	cfg            *config.Config
	guardSystem    *GuardSystem
	permissionCtrl *PermissionController
	executionCtrl  *ExecutionController
	evolver        *SchemaEvolver
	auditLogger    *AuditLogger
	cache          *QueryCache
	semTopology    *SemanticTopology
	closed         bool
}

func NewText2SQLSkill(cfg *config.Config, db *sql.DB) (interfaces.Skill, error) {
	permCtrl := NewPermissionController(cfg)
	execCtrl := NewExecutionController(cfg)
	guardSystem := NewGuardSystem(cfg, permCtrl, execCtrl)
	evolver := NewSchemaEvolver(cfg)
	auditLogger := NewAuditLogger(cfg)
	cache := NewQueryCache(cfg)
	semTopology := NewSemanticTopology()

	return &Text2SQLSkill{
		db:             db,
		cfg:            cfg,
		guardSystem:    guardSystem,
		permissionCtrl: permCtrl,
		executionCtrl:  execCtrl,
		evolver:        evolver,
		auditLogger:    auditLogger,
		cache:          cache,
		semTopology:    semTopology,
	}, nil
}

func (s *Text2SQLSkill) CapabilityID() string {
	return s.cfg.App.Name + "-" + s.cfg.App.Version
}

func (s *Text2SQLSkill) Execute(ctx context.Context, input string) (interfaces.SkillResult, error) {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return interfaces.SkillResult{}, fmt.Errorf("skill is closed")
	}
	s.mu.Unlock()

	queryID := utils.GenerateQueryID()
	startTime := time.Now()

	if s.cfg.Audit.Enabled {
		s.auditLogger.LogEvent(queryID, "execution_start", map[string]interface{}{
			"input": input,
		})
	}

	defer func() {
		duration := time.Since(startTime)
		if s.cfg.Audit.Enabled {
			if s.cfg.Performance.AsyncProcessing {
				go s.auditLogger.LogEvent(queryID, "execution_end", map[string]interface{}{
					"duration_ms": duration.Milliseconds(),
				})
			} else {
				s.auditLogger.LogEvent(queryID, "execution_end", map[string]interface{}{
					"duration_ms": duration.Milliseconds(),
				})
			}
		}
	}()

	// Check cache first
	if s.cfg.Cache.Enabled {
		if result, found := s.cache.Get(input); found {
			if s.cfg.Audit.Enabled {
				s.auditLogger.LogEvent(queryID, "cache_hit", map[string]interface{}{
					"input": input,
				})
			}
			return result, nil
		}
	}

	// Five layer guard check
	if allowed, reason := s.guardSystem.CheckAllGuards(ctx, input); !allowed {
		result := interfaces.SkillResult{
			QueryID:   queryID,
			Meta:      []byte(reason),
			Timestamp: time.Now(),
			Status:    "rejected",
		}

		if s.cfg.Audit.Enabled {
			s.auditLogger.LogEvent(queryID, "rejected", map[string]interface{}{
				"input":  input,
				"reason": reason,
			})
		}

		return result, nil
	}

	// Build semantic topology
	topology := s.semTopology.BuildTopology([]byte(input))
	if topology == nil {
		result := interfaces.SkillResult{
			QueryID:   queryID,
			Meta:      []byte("topology_build_failed"),
			Timestamp: time.Now(),
			Status:    "error",
		}

		if s.cfg.Audit.Enabled {
			s.auditLogger.LogEvent(queryID, "topology_error", map[string]interface{}{
				"input": input,
			})
		}

		return result, nil
	}

	// Generate query template
	fingerprint := s.semTopology.GenerateTopologyFingerprint(topology)
	template := s.evolver.GetQueryTemplate(fingerprint)

	// Execute with isolation
	execCtx, cancel := s.executionCtrl.GetExecutionContext(ctx)
	defer cancel()

	rows, err := s.executeQueryWithIsolation(execCtx, template, input)
	if err != nil {
		result := interfaces.SkillResult{
			QueryID:   queryID,
			Meta:      []byte("execution_failed: " + err.Error()),
			Timestamp: time.Now(),
			Status:    "error",
		}

		if s.cfg.Audit.Enabled {
			s.auditLogger.LogEvent(queryID, "execution_error", map[string]interface{}{
				"input":   input,
				"error":   err.Error(),
				"timeout": s.cfg.Execution.Timeout.Total,
			})
		}

		return result, nil
	}

	// Process results
	resultData := s.processResultRows(rows)
	// 使用安全配置中的资源限制
	maxRows := s.cfg.Security.ResourceLimits.MaxRows
	if len(resultData) > maxRows {
		resultData = resultData[:maxRows]
	}

	// Generate encrypted result
	compress := s.cfg.Performance.Compression.Enabled
	encryptedResult := utils.EncryptResult(resultData, compress)

	// Create result
	result := interfaces.SkillResult{
		QueryID:   queryID,
		Result:    encryptedResult,
		Meta:      s.generateMetadata(input, template, len(resultData)),
		Timestamp: time.Now(),
		Status:    "success",
	}

	// Cache result
	if s.cfg.Cache.Enabled {
		s.cache.Set(input, result)
	}

	// Audit success
	if s.cfg.Audit.Enabled {
		s.auditLogger.LogEvent(queryID, "success", map[string]interface{}{
			"input":       input,
			"template":    template,
			"row_count":   len(resultData),
			"duration_ms": time.Since(startTime).Milliseconds(),
		})
	}

	return result, nil
}

func (s *Text2SQLSkill) executeQueryWithIsolation(ctx context.Context, template string, input string) (*sql.Rows, error) {
	switch s.executionCtrl.GetIsolationLevel() {
	case "full":
		return s.executeQueryWithFullIsolation(ctx, template, input)
	case "basic":
		return s.executeQueryWithBasicIsolation(ctx, template, input)
	default:
		return s.db.QueryContext(ctx, template)
	}
}

func (s *Text2SQLSkill) executeQueryWithFullIsolation(ctx context.Context, template string, input string) (*sql.Rows, error) {
	resultChan := make(chan struct {
		rows *sql.Rows
		err  error
	}, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				resultChan <- struct {
					rows *sql.Rows
					err  error
				}{nil, fmt.Errorf("execution panic: %v", r)}
			}
		}()

		rows, err := s.db.QueryContext(ctx, template)
		resultChan <- struct {
			rows *sql.Rows
			err  error
		}{rows, err}
	}()

	// 解析超时时间
	timeout, err := time.ParseDuration(s.cfg.Execution.Timeout.Total)
	if err != nil {
		timeout = 10 * time.Second // 默认值
	}

	select {
	case result := <-resultChan:
		return result.rows, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(timeout):
		return nil, fmt.Errorf("execution timeout after %v", timeout)
	}
}

func (s *Text2SQLSkill) executeQueryWithBasicIsolation(ctx context.Context, template string, input string) (*sql.Rows, error) {
	rows, err := s.db.QueryContext(ctx, template)
	if err != nil {
		return nil, err
	}

	// Row count limit - 使用安全配置中的资源限制
	maxRows := s.cfg.Security.ResourceLimits.MaxRows
	count := 0
	for rows.Next() {
		count++
		if count >= maxRows {
			break
		}
	}

	return rows, nil
}

func (s *Text2SQLSkill) processResultRows(rows *sql.Rows) []map[string]interface{} {
	defer rows.Close()

	columns, _ := rows.Columns()
	types, _ := rows.ColumnTypes()

	var results []map[string]interface{}

	// 使用安全配置中的资源限制
	maxRows := s.cfg.Security.ResourceLimits.MaxRows

	for rows.Next() && len(results) < maxRows {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i, typ := range types {
			switch typ.DatabaseTypeName() {
			case "INT", "BIGINT", "TINYINT", "SMALLINT", "MEDIUMINT":
				values[i] = new(int64)
			case "DECIMAL", "FLOAT", "DOUBLE", "REAL":
				values[i] = new(float64)
			case "VARCHAR", "CHAR", "TEXT", "MEDIUMTEXT", "LONGTEXT":
				values[i] = new(string)
			default:
				values[i] = new(string)
			}
			valuePtrs[i] = values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			switch v := values[i].(type) {
			case *int64:
				row[col] = *v
			case *float64:
				row[col] = *v
			case *string:
				row[col] = *v
			}
		}
		results = append(results, row)
	}

	return results
}

func (s *Text2SQLSkill) generateMetadata(input string, template string, rowCount int) []byte {
	metadata := map[string]interface{}{
		"input_length":  len(input),
		"template_used": template,
		"row_count":     rowCount,
		"timestamp":     time.Now().UTC().Format(time.RFC3339Nano),
	}

	data, _ := json.Marshal(metadata)
	return data
}

func (s *Text2SQLSkill) SafeShutdown() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	if s.auditLogger != nil {
		s.auditLogger.Close()
	}

	if s.db != nil {
		s.db.Close()
	}

	s.closed = true
	return nil
}
