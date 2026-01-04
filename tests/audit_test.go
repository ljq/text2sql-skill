package tests

import (
	"testing"

	"text2sql-skill/config"
	"text2sql-skill/core"
)

func TestAuditLogger(t *testing.T) {
	// 使用内存中的模拟审计日志器，避免文件系统操作
	cfg := config.DefaultConfig()
	cfg.Audit.Storage.Type = "console" // 使用控制台输出代替文件
	cfg.Audit.Enabled = true
	cfg.Audit.Level = "detailed"

	logger := core.NewAuditLogger(cfg)
	defer logger.Close()

	// Test logging
	queryID := "test_query_123"
	logger.LogEvent(queryID, "test_event", map[string]interface{}{
		"test": "value",
		"num":  42,
	})

	// 验证日志器已初始化
	if logger == nil {
		t.Error("Audit logger failed to initialize")
	}
}
