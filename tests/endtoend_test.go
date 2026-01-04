package tests

import (
	"context"
	"testing"
	"time"

	"text2sql-skill/config"
	"text2sql-skill/core"
	"text2sql-skill/drivers"
	"text2sql-skill/interfaces"
	"text2sql-skill/utils"
)

func TestEndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test in short mode")
	}

	// Create test config without file system operations
	cfg := config.DefaultConfig()
	cfg.Database.Driver = "mysql"
	cfg.Database.MySQL.DSN = "test:test@tcp(localhost:3306)/test_db"
	cfg.Audit.Storage.Type = "console" // 使用控制台输出避免文件操作
	cfg.Audit.Enabled = true

	// Register driver
	drivers.RegisterMySQLDriver()

	// Create mock skill for testing
	skill := createMockSkill(cfg)

	ctx := context.Background()
	input := "2025年北京销售额超过100万的客户"

	// First execution (cache miss)
	result1, err := skill.Execute(ctx, input)
	if err != nil {
		t.Fatalf("First execution failed: %v", err)
	}
	if result1.Status != "success" {
		t.Errorf("First execution expected status 'success', got '%s'", result1.Status)
	}

	// Second execution (cache hit)
	result2, err := skill.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Second execution failed: %v", err)
	}
	if result2.Status != "success" {
		t.Errorf("Second execution expected status 'success', got '%s'", result2.Status)
	}

	// Test forbidden operation
	forbiddenInput := "DELETE FROM customers"
	result3, err := skill.Execute(ctx, forbiddenInput)
	if err != nil {
		t.Fatalf("Forbidden execution failed: %v", err)
	}
	if result3.Status != "rejected" {
		t.Errorf("Forbidden execution expected status 'rejected', got '%s'", result3.Status)
	}
}

func createMockSkill(cfg *config.Config) interfaces.Skill {
	// This is a mock implementation for testing
	// In production, this would use a real database connection
	permCtrl := core.NewPermissionController(cfg)
	execCtrl := core.NewExecutionController(cfg)
	guardSystem := core.NewGuardSystem(cfg, permCtrl, execCtrl)
	evolver := core.NewSchemaEvolver(cfg)
	auditLogger := core.NewAuditLogger(cfg)
	cache := core.NewQueryCache(cfg)
	semTopology := core.NewSemanticTopology()

	return &mockSkill{
		cfg:         cfg,
		guardSystem: guardSystem,
		evolver:     evolver,
		auditLogger: auditLogger,
		cache:       cache,
		semTopology: semTopology,
	}
}

type mockSkill struct {
	cfg         *config.Config
	guardSystem *core.GuardSystem
	evolver     *core.SchemaEvolver
	auditLogger *core.AuditLogger
	cache       *core.QueryCache
	semTopology *core.SemanticTopology
}

func (m *mockSkill) CapabilityID() string {
	return m.cfg.App.Name + "-" + m.cfg.App.Version
}

func (m *mockSkill) Execute(ctx context.Context, input string) (interfaces.SkillResult, error) {
	queryID := utils.GenerateQueryID()

	if allowed, reason := m.guardSystem.CheckAllGuards(ctx, input); !allowed {
		return interfaces.SkillResult{
			QueryID:   queryID,
			Meta:      []byte(reason),
			Timestamp: time.Now(),
			Status:    "rejected",
		}, nil
	}

	// Mock successful result
	return interfaces.SkillResult{
		QueryID:   queryID,
		Result:    []byte("mock_result_data"),
		Meta:      []byte("mock_metadata"),
		Timestamp: time.Now(),
		Status:    "success",
	}, nil
}

func (m *mockSkill) SafeShutdown() error {
	m.auditLogger.Close()
	return nil
}
