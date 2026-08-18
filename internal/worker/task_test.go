package worker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskPayload(t *testing.T) {
	payload := NewTaskPayload(TaskTypePriceCheck, "product-123", "")

	assert.NotEmpty(t, payload.TaskID)
	assert.Equal(t, TaskTypePriceCheck, payload.Type)
	assert.Equal(t, "product-123", payload.ProductID)
	assert.False(t, payload.CreatedAt.IsZero())
}

func TestTaskPayload_ToJSON(t *testing.T) {
	payload := NewTaskPayload(TaskTypeCompetitorSync, "", "search query")
	payload.TaskID = "fixed-id"

	jsonStr, err := payload.ToJSON()
	assert.NoError(t, err)
	assert.Contains(t, jsonStr, "fixed-id")
	assert.Contains(t, jsonStr, "competitor_sync")
	assert.Contains(t, jsonStr, "search query")
}

func TestTaskPayload_ToJSON_EmptyProductIDAndQuery(t *testing.T) {
	// Ensure ToJSON works with empty optional fields
	payload := NewTaskPayload(TaskTypePriceCheck, "", "")
	payload.TaskID = "test-id"

	jsonStr, err := payload.ToJSON()
	assert.NoError(t, err)
	assert.Contains(t, jsonStr, "test-id")
	assert.Contains(t, jsonStr, "price_check")
}

func TestParseTaskPayload(t *testing.T) {
	original := NewTaskPayload(TaskTypeTrendAnalysis, "prod-456", "analyze trends")
	original.TaskID = "test-id"

	jsonStr, err := original.ToJSON()
	assert.NoError(t, err)

	parsed, err := ParseTaskPayload(jsonStr)
	assert.NoError(t, err)
	assert.Equal(t, original.TaskID, parsed.TaskID)
	assert.Equal(t, original.Type, parsed.Type)
	assert.Equal(t, original.ProductID, parsed.ProductID)
	assert.Equal(t, original.Query, parsed.Query)
}

func TestParseTaskPayload_Invalid(t *testing.T) {
	_, err := ParseTaskPayload("invalid json")
	assert.Error(t, err)
}

func TestTaskType_Constants(t *testing.T) {
	assert.Equal(t, TaskType("price_check"), TaskTypePriceCheck)
	assert.Equal(t, TaskType("competitor_sync"), TaskTypeCompetitorSync)
	assert.Equal(t, TaskType("trend_analysis"), TaskTypeTrendAnalysis)
	assert.Equal(t, TaskType("supply_alert"), TaskTypeSupplyAlert)
}

func TestDefaultRetryPolicy(t *testing.T) {
	policy := DefaultRetryPolicy
	assert.Equal(t, 3, policy.MaxRetry)
	assert.Equal(t, 10*time.Second, policy.Initial)
	assert.Equal(t, 10*time.Minute, policy.MaxDelay)
	assert.Equal(t, 2.0, policy.Multiplier)
}
