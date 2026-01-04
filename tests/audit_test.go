package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"text2sql-skill/config"
	"text2sql-skill/core"
)

func TestAuditLogger(t *testing.T) {
	// Setup test directory
	testDir := filepath.Join(os.TempDir(), "audit_test_"+time.Now().Format("20060102150405"))
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	cfg := config.DefaultConfig()
	cfg.Audit.Storage.Type = "file"
	cfg.Audit.Storage.Path = testDir
	cfg.Audit.Storage.Rotation.MaxSizeMB = 10
	cfg.Audit.Storage.Rotation.MaxAgeDays = 1
	cfg.Audit.Storage.Rotation.MaxBackups = 5
	cfg.Audit.Storage.Rotation.Compress = false

	logger := core.NewAuditLogger(cfg)
	defer logger.Close()

	// Test logging
	queryID := "test_query_123"
	logger.LogEvent(queryID, "test_event", map[string]interface{}{
		"test": "value",
		"num":  42,
	})

	// Check if file was created
	day := time.Now().UTC().Format("2006-01-02")
	filename := filepath.Join(testDir, "audit_"+day+".log")

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("Audit log file was not created")
	}
}
