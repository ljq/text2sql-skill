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
	"strings"
	"time"

	"text2sql-skill/config"
)

type GuardLevel int

const (
	GuardL1_SemanticSafety GuardLevel = iota + 1
	GuardL2_OperationPermission
	GuardL3_KeywordFilter
	GuardL4_ResourceControl
	GuardL5_ExecutionSafety
)

type GuardSystem struct {
	cfg            *config.Config
	permissionCtrl *PermissionController
	executionCtrl  *ExecutionController
}

func NewGuardSystem(cfg *config.Config, permCtrl *PermissionController, execCtrl *ExecutionController) *GuardSystem {
	return &GuardSystem{
		cfg:            cfg,
		permissionCtrl: permCtrl,
		executionCtrl:  execCtrl,
	}
}

func (g *GuardSystem) CheckAllGuards(ctx context.Context, input string) (bool, string) {
	inputBytes := []byte(input)

	// L1: Semantic Safety
	if !g.permissionCtrl.CheckSemanticSafety(inputBytes) {
		return false, "L1: semantic safety violation - entropy or ratio out of configured range"
	}

	// L2: Operation Permission
	operation := g.detectOperationType(inputBytes)
	if !g.permissionCtrl.CheckOperationPermission(operation) {
		return false, "L2: operation not allowed in current execution mode"
	}

	// L3: Keyword Filter
	if keyword := g.permissionCtrl.CheckForbiddenKeywords(inputBytes); keyword != "" {
		return false, "L3: forbidden keyword detected: " + keyword
	}

	// L4: Resource Control
	if !g.checkResourceLimits(inputBytes) {
		return false, "L4: resource limits exceeded"
	}

	// L5: Execution Safety
	if err := g.checkExecutionSafety(ctx); err != nil {
		return false, "L5: " + err.Error()
	}

	return true, ""
}

func (g *GuardSystem) detectOperationType(input []byte) string {
	lowerInput := strings.ToLower(string(input))

	for _, op := range g.cfg.Security.AllowedOperations {
		if strings.Contains(lowerInput, strings.ToLower(op)) {
			return op
		}
	}

	return "SELECT"
}

func (g *GuardSystem) checkResourceLimits(input []byte) bool {
	inputSize := len(input)
	estimatedRows := inputSize / 100
	estimatedMemoryMB := float64(inputSize) / (1024 * 1024)

	return g.executionCtrl.CheckResourceLimits(inputSize, estimatedRows, estimatedMemoryMB)
}

func (g *GuardSystem) checkExecutionSafety(ctx context.Context) error {
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		// 解析超时时间，使用总超时的一半作为最小要求
		timeout, err := time.ParseDuration(g.cfg.Execution.Timeout.Total)
		if err != nil {
			timeout = 10 * time.Second // 默认值
		}
		minRequired := timeout / 2
		if remaining < minRequired {
			return context.DeadlineExceeded
		}
	}
	return nil
}
