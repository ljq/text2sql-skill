package core

import (
	"sync"
	"time"

	"text2sql-skill/config"
	"text2sql-skill/interfaces"
)

type QueryCache struct {
	sync.RWMutex
	cfg    *config.Config
	cache  map[string]interfaces.SkillResult
	ttl    time.Duration
	maxAge time.Time
}

func NewQueryCache(cfg *config.Config) *QueryCache {
	cache := &QueryCache{
		cfg:   cfg,
		cache: make(map[string]interfaces.SkillResult),
		ttl:   5 * time.Minute,
	}

	// 解析 TTL
	if cfg.Cache.Enabled && cfg.Cache.TTL != "" {
		if ttl, err := time.ParseDuration(cfg.Cache.TTL); err == nil {
			cache.ttl = ttl
		}
	}

	if cfg.Cache.Enabled && cfg.Cache.Size > 0 {
		go cache.cleanupLoop()
	}

	return cache
}

func (c *QueryCache) Get(input string) (interfaces.SkillResult, bool) {
	c.RLock()
	defer c.RUnlock()

	result, found := c.cache[input]
	if found && time.Now().Before(result.Timestamp.Add(c.ttl)) {
		return result, true
	}

	return interfaces.SkillResult{}, false
}

func (c *QueryCache) Set(input string, result interfaces.SkillResult) {
	c.Lock()
	defer c.Unlock()

	if !c.cfg.Cache.Enabled {
		return
	}

	if len(c.cache) >= c.cfg.Cache.Size {
		c.evictOldest()
	}

	c.cache[input] = result
}

func (c *QueryCache) evictOldest() {
	oldest := time.Now()
	var oldestKey string

	for key, result := range c.cache {
		if result.Timestamp.Before(oldest) {
			oldest = result.Timestamp
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(c.cache, oldestKey)
	}
}

func (c *QueryCache) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.Lock()
		now := time.Now()
		for key, result := range c.cache {
			if now.After(result.Timestamp.Add(c.ttl)) {
				delete(c.cache, key)
			}
		}
		c.Unlock()

		if len(c.cache) == 0 {
			break
		}
	}
}
