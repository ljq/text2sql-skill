package interfaces

import (
	"context"
	"time"
)

type Skill interface {
	CapabilityID() string
	Execute(ctx context.Context, input string) (SkillResult, error)
	SafeShutdown() error
}

type SkillResult struct {
	QueryID   string    `json:"query_id"`
	Result    []byte    `json:"result"`
	Meta      []byte    `json:"meta"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}
