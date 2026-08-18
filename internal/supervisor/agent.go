package supervisor

import (
	"context"
	"time"
)

// TaskType represents the type of task.
type TaskType string

const (
	TaskTypeMarketTrend    TaskType = "market_trend"
	TaskTypeCompetitor     TaskType = "competitor"
	TaskTypePricing        TaskType = "pricing"
	TaskTypeSupplyChain    TaskType = "supply_chain"
)

// Task represents a subtask decomposed from user query.
type Task struct {
	ID          string    `json:"id"`
	Type        TaskType  `json:"type"`
	Description string    `json:"description"`
	Query       string    `json:"query"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	Result      any       `json:"result,omitempty"`
}

// Agent defines the interface for agents that execute tasks.
type Agent interface {
	Execute(ctx context.Context, task *Task) (any, error)
	Name() string
}
