package worker

import (
	"context"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
)

func TestNewTaskHandler(t *testing.T) {
	logger := newTestLogger()
	handler := NewHandler(logger, nil, nil)
	taskHandler := NewTaskHandler(logger, handler)

	assert.NotNil(t, taskHandler)
	assert.Equal(t, logger, taskHandler.logger)
	assert.Equal(t, handler, taskHandler.inner)
}

func TestTaskHandler_ProcessTask(t *testing.T) {
	logger := newTestLogger()
	handler := NewHandler(logger, nil, nil)
	taskHandler := NewTaskHandler(logger, handler)

	tests := []struct {
		name        string
		taskType    TaskType
		productID   string
		expectError bool
	}{
		{
			name:        "price check",
			taskType:    TaskTypePriceCheck,
			productID:   "test-prod-1",
			expectError: false,
		},
		{
			name:        "competitor sync",
			taskType:    TaskTypeCompetitorSync,
			expectError: false,
		},
		{
			name:        "trend analysis",
			taskType:    TaskTypeTrendAnalysis,
			expectError: true, // nil LLM returns error
		},
		{
			name:        "supply alert",
			taskType:    TaskTypeSupplyAlert,
			productID:   "test-prod-2",
			expectError: true, // nil LLM returns error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := NewTaskPayload(tt.taskType, tt.productID, "test query")
			payloadJSON, _ := payload.ToJSON()

			task := asynq.NewTask(string(tt.taskType), []byte(payloadJSON))
			err := taskHandler.ProcessTask(context.Background(), task)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskHandler_ProcessTask_InvalidPayload(t *testing.T) {
	logger := newTestLogger()
	handler := NewHandler(logger, nil, nil)
	taskHandler := NewTaskHandler(logger, handler)

	task := asynq.NewTask("price_check", []byte("invalid json"))
	err := taskHandler.ProcessTask(context.Background(), task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}
