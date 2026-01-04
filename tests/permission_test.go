package tests

import (
	"context"
	"testing"

	"text2sql-skill/config"
	"text2sql-skill/core"
)

func TestPermissionController(t *testing.T) {
	cfg := config.DefaultConfig()
	permCtrl := core.NewPermissionController(cfg)

	tests := []struct {
		input   string
		allowed bool
		reason  string
	}{
		{"SELECT customers FROM sales WHERE region = '北京'", true, ""},
		{"DELETE FROM customers WHERE id = 1", false, "L2: operation not allowed in current execution mode"},
		{"DROP TABLE sales", false, "L3: forbidden keyword detected: DROP"},
		{"UPDATE customers SET name = 'test'", false, "L3: forbidden keyword detected: UPDATE"},
		{"SELECT * FROM sales -- DROP TABLE customers", false, "L3: forbidden keyword detected: DROP"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			allowed, _ := core.NewGuardSystem(cfg, permCtrl, core.NewExecutionController(cfg)).CheckAllGuards(context.Background(), tt.input)
			if allowed != tt.allowed {
				t.Errorf("Input: %q, Expected: %v, Got: %v", tt.input, tt.allowed, allowed)
			}
		})
	}
}

func TestSemanticSafety(t *testing.T) {
	cfg := config.DefaultConfig()
	permCtrl := core.NewPermissionController(cfg)

	// Valid Chinese query
	validInput := []byte("2025年北京销售额超过100万的客户")
	if !permCtrl.CheckSemanticSafety(validInput) {
		t.Error("Valid Chinese query should pass semantic safety check")
	}

	// Invalid input (too short)
	invalidInput := []byte("a")
	if permCtrl.CheckSemanticSafety(invalidInput) {
		t.Error("Invalid short input should fail semantic safety check")
	}

	// English-only input (low non-ASCII ratio)
	englishInput := []byte("SELECT customers FROM sales WHERE region = 'Beijing'")
	if permCtrl.CheckSemanticSafety(englishInput) {
		t.Error("English-only input should fail semantic safety check in strict mode")
	}
}
