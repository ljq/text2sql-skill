package core

import (
	"sync"

	"text2sql-skill/config"
)

type SchemaEvolver struct {
	sync.RWMutex
	cfg       *config.Config
	mappings  map[string]uint32
	counters  map[uint32]int
	patternID uint32
}

func NewSchemaEvolver(cfg *config.Config) *SchemaEvolver {
	return &SchemaEvolver{
		cfg:      cfg,
		mappings: make(map[string]uint32),
		counters: make(map[uint32]int),
	}
}

func (e *SchemaEvolver) GetQueryTemplate(fingerprint []byte) string {
	e.RLock()
	defer e.RUnlock()

	fingerprintStr := string(fingerprint)
	if schemaID, exists := e.mappings[fingerprintStr]; exists {
		return e.generateTemplate(schemaID)
	}
	return "SELECT * FROM data WHERE 1=1"
}

func (e *SchemaEvolver) RegisterNewPattern(fingerprint []byte) {
	e.Lock()
	defer e.Unlock()

	fingerprintStr := string(fingerprint)
	if _, exists := e.mappings[fingerprintStr]; exists {
		return
	}

	// 新配置中移除了 Evolution 部分，使用固定值
	maxPatterns := 5000 // 默认值
	if len(e.mappings) >= maxPatterns {
		e.evictOldestPattern()
	}

	e.patternID++
	e.mappings[fingerprintStr] = e.patternID
	e.counters[e.patternID] = 1
}

func (e *SchemaEvolver) evictOldestPattern() {
	minCount := int(^uint(0) >> 1)
	var oldestID uint32

	for id, count := range e.counters {
		if count < minCount {
			minCount = count
			oldestID = id
		}
	}

	for fingerprint, id := range e.mappings {
		if id == oldestID {
			delete(e.mappings, fingerprint)
			delete(e.counters, id)
			break
		}
	}
}

func (e *SchemaEvolver) generateTemplate(schemaID uint32) string {
	switch {
	case schemaID%3 == 0:
		return "SELECT name FROM customers c JOIN sales s ON c.id = s.customer_id WHERE s.region = ? AND s.year = ? AND s.amount > ?"
	case schemaID%3 == 1:
		return "SELECT region, SUM(amount) as total FROM sales WHERE year = ? GROUP BY region"
	default:
		return "SELECT name, amount FROM customers c JOIN sales s ON c.id = s.customer_id WHERE s.year = ? ORDER BY s.amount DESC LIMIT ?"
	}
}
