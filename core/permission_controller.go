package core

import (
	"math"
	"strings"
	"unicode"
	"unicode/utf8"

	"text2sql-skill/config"
)

type PermissionController struct {
	cfg *config.Config
}

func NewPermissionController(cfg *config.Config) *PermissionController {
	return &PermissionController{cfg: cfg}
}

func (p *PermissionController) CheckOperationPermission(operation string) bool {
	switch p.cfg.Security.Mode {
	case "read_only":
		return strings.EqualFold(operation, "SELECT")
	case "read_write":
		return true
	default:
		// 对于其他模式，检查是否在允许的操作列表中
		for _, allowed := range p.cfg.Security.AllowedOperations {
			if strings.EqualFold(allowed, operation) {
				return true
			}
		}
		return false
	}
}

func (p *PermissionController) CheckSemanticSafety(input []byte) bool {
	entropy := p.calculateEntropy(input)
	// nonASCIIRatio := p.calculateNonASCIIRatio(input) // 新配置中移除了此检查

	validation := p.cfg.Security.InputValidation
	return entropy >= validation.MinEntropy &&
		entropy <= validation.MaxEntropy
}

func (p *PermissionController) CheckTopologyBalance(balance float32) bool {
	// 新配置中移除了 topology_balance，暂时返回 true
	// 如果需要此功能，可以在后续版本中添加
	return true
}

func (p *PermissionController) CheckForbiddenKeywords(input []byte) string {
	lowerInput := strings.ToLower(string(input))
	for _, keyword := range p.cfg.Security.ForbiddenKeywords {
		if strings.Contains(lowerInput, strings.ToLower(keyword)) {
			return keyword
		}
	}
	return ""
}

func (p *PermissionController) calculateEntropy(input []byte) float32 {
	runeCount := make(map[rune]int)
	total := 0

	for i := 0; i < len(input); {
		r, size := utf8.DecodeRune(input[i:])
		if !unicode.IsSpace(r) && !unicode.IsPunct(r) {
			runeCount[r]++
			total++
		}
		i += size
	}

	if total == 0 {
		return 0
	}

	var entropy float32
	for _, count := range runeCount {
		p := float32(count) / float32(total)
		if p > 0 {
			entropy -= p * float32(log2(float64(p)))
		}
	}

	return entropy
}

func (p *PermissionController) calculateNonASCIIRatio(input []byte) float32 {
	nonASCII := 0
	total := 0

	for i := 0; i < len(input); {
		r, size := utf8.DecodeRune(input[i:])
		if !unicode.IsSpace(r) && !unicode.IsPunct(r) {
			total++
			if r > 127 {
				nonASCII++
			}
		}
		i += size
	}

	if total == 0 {
		return 0
	}

	return float32(nonASCII) / float32(total)
}

func log2(x float64) float64 {
	return math.Log(x) / math.Log(2)
}
