package tests

import (
	"context"
	"testing"
	"time"

	"text2sql-skill/config"
	"text2sql-skill/core"
)

func TestGuardSystem(t *testing.T) {
	cfg := config.DefaultConfig()
	permCtrl := core.NewPermissionController(cfg)
	execCtrl := core.NewExecutionController(cfg)
	guardSystem := core.NewGuardSystem(cfg, permCtrl, execCtrl)

	tests := []struct {
		name    string
		input   string
		ctx     context.Context
		allowed bool
		reason  string
	}{
		{
			name:    "valid query",
			input:   "2025年北京销售额超过100万的客户",
			ctx:     context.Background(),
			allowed: true,
		},
		{
			name:    "forbidden keyword",
			input:   "DELETE FROM customers",
			ctx:     context.Background(),
			allowed: false,
			reason:  "L3: forbidden keyword detected: DELETE",
		},
		{
			name:    "timeout context",
			input:   "long running query",
			ctx:     createTimeoutContext(100),
			allowed: false,
			reason:  "L5: context deadline exceeded",
		},
		{
			name:    "large input",
			input:   generateLargeInput(20000),
			ctx:     context.Background(),
			allowed: false,
			reason:  "L4: resource limits exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed, reason := guardSystem.CheckAllGuards(tt.ctx, tt.input)
			if allowed != tt.allowed {
				t.Errorf("Expected allowed=%v, got allowed=%v", tt.allowed, allowed)
			}
			if tt.reason != "" && reason != tt.reason {
				t.Errorf("Expected reason=%q, got reason=%q", tt.reason, reason)
			}
		})
	}
}

func createTimeoutContext(ms int) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ms)*time.Millisecond)
	defer cancel()
	return ctx
}

func generateLargeInput(size int) string {
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = 'a'
	}
	return string(buf)
}
