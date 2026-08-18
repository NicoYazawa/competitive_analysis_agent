package supervisor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTask_Structure(t *testing.T) {
	task := &Task{
		ID:          "task-123",
		Type:        TaskTypeMarketTrend,
		Description: "Test task",
		Query:       "test query",
		Priority:    1,
		CreatedAt:   time.Now(),
	}

	assert.Equal(t, "task-123", task.ID)
	assert.Equal(t, TaskTypeMarketTrend, task.Type)
	assert.Equal(t, "Test task", task.Description)
	assert.Equal(t, "test query", task.Query)
	assert.Equal(t, 1, task.Priority)
	assert.False(t, task.CreatedAt.IsZero())
	assert.Nil(t, task.Result)
}

func TestTaskType_Constants(t *testing.T) {
	assert.Equal(t, TaskType("market_trend"), TaskTypeMarketTrend)
	assert.Equal(t, TaskType("competitor"), TaskTypeCompetitor)
	assert.Equal(t, TaskType("pricing"), TaskTypePricing)
	assert.Equal(t, TaskType("supply_chain"), TaskTypeSupplyChain)
}

func TestScheduler_containsAny(t *testing.T) {
	tests := []struct {
		s        string
		keywords []string
		expected bool
	}{
		{"hello world", []string{"hello", "world"}, true},
		{"hello world", []string{"foo", "bar"}, false},
		{"", []string{"hello"}, false},
		{"HELLO", []string{"hello"}, false}, // case sensitive
	}

	for _, tt := range tests {
		result := containsAny(tt.s, tt.keywords)
		assert.Equal(t, tt.expected, result)
	}
}
