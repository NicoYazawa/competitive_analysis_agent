package supervisor

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScheduler_DecomposeQuery(t *testing.T) {
	s := NewScheduler()

	tests := []struct {
		name         string
		query        string
		expectedTypes []TaskType
	}{
		{
			name:          "market trend query",
			query:         "analyze market trend for smartphones",
			expectedTypes: []TaskType{TaskTypeMarketTrend},
		},
		{
			name:          "competitor query",
			query:         "compare competitors on Amazon",
			expectedTypes: []TaskType{TaskTypeCompetitor},
		},
		{
			name:          "pricing query",
			query:         "recommend pricing strategy",
			expectedTypes: []TaskType{TaskTypePricing},
		},
		{
			name:          "supply chain query",
			query:         "analyze supply chain risks",
			expectedTypes: []TaskType{TaskTypeSupplyChain},
		},
		{
			name:          "multi-task query",
			query:         "market trend and competitor analysis for electronics",
			expectedTypes: []TaskType{TaskTypeMarketTrend, TaskTypeCompetitor},
		},
		{
			name:          "Chinese keywords",
			query:         "分析市场趋势和竞品",
			expectedTypes: []TaskType{TaskTypeMarketTrend, TaskTypeCompetitor},
		},
		{
			name:          "unknown query defaults to market trend",
			query:         "something random",
			expectedTypes: []TaskType{TaskTypeMarketTrend},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := s.DecomposeQuery(context.Background(), tt.query)
			require.NoError(t, err)
			assert.Len(t, tasks, len(tt.expectedTypes))

			for i, expectedType := range tt.expectedTypes {
				assert.Equal(t, expectedType, tasks[i].Type)
			}
		})
	}
}

func TestScheduler_RegisterAgent(t *testing.T) {
	s := NewScheduler()

	mockAgent := &mockAgent{name: "test-agent"}
	s.RegisterAgent(TaskTypeMarketTrend, mockAgent)

	_, err := s.DecomposeQuery(context.Background(), "test query")
	require.NoError(t, err)

	result, err := s.ScheduleAndExecute(context.Background(), "test query")
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestScheduler_ScheduleAndExecute(t *testing.T) {
	s := NewScheduler()

	mockAgent := &mockAgent{
		name: "market-agent",
		result: &MarketTrendResult{
			Trend:         "Bullish trend",
			Opportunities: []string{"Opportunity 1"},
			DemandSignal:  "Strong",
		},
	}
	s.RegisterAgent(TaskTypeMarketTrend, mockAgent)

	result, err := s.ScheduleAndExecute(context.Background(), "analyze market trend for phones")

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Bullish trend", result[0].Result.(*MarketTrendResult).Trend)
}

func TestScheduler_ScheduleAndExecute_Error(t *testing.T) {
	s := NewScheduler()

	mockAgent := &mockAgent{
		name: "failing-agent",
		err:  errors.New("agent failed"),
	}
	s.RegisterAgent(TaskTypeMarketTrend, mockAgent)

	result, err := s.ScheduleAndExecute(context.Background(), "analyze market trend")

	require.NoError(t, err)
	assert.Len(t, result, 1)
	// Result is a wrapped error, not a string
	errResult, ok := result[0].Result.(error)
	require.True(t, ok)
	assert.Contains(t, errResult.Error(), "agent failed")
}

type mockAgent struct {
	name   string
	result any
	err    error
}

func (a *mockAgent) Name() string { return a.name }

func (a *mockAgent) Execute(ctx context.Context, task *Task) (any, error) {
	if a.err != nil {
		return nil, a.err
	}
	return a.result, nil
}
