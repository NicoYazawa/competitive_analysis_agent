package agents

import (
	"context"
	"errors"
	"testing"

	"competitive-analysis-agent/internal/llm"
	"competitive-analysis-agent/internal/supervisor"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarketAgent_Execute(t *testing.T) {
	mockLLM := &llm.MockLLM{
		RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
			return &llm.ChatResponse{Content: "Market trend is bullish"}, nil
		},
	}

	agent := NewMarketAgent(mockLLM)
	task := &supervisor.Task{
		ID:    "task-1",
		Type:  supervisor.TaskTypeMarketTrend,
		Query: "analyze smartphone market",
	}

	result, err := agent.Execute(context.Background(), task)

	require.NoError(t, err)
	marketResult, ok := result.(*supervisor.MarketTrendResult)
	require.True(t, ok)
	assert.NotEmpty(t, marketResult.Trend)
	assert.NotEmpty(t, marketResult.DemandSignal)
}

func TestMarketAgent_Execute_LLMError(t *testing.T) {
	mockLLM := &llm.MockLLM{
		ErrFn: func(ctx context.Context, req *llm.ChatRequest) error {
			return errors.New("LLM error")
		},
	}

	agent := NewMarketAgent(mockLLM)
	task := &supervisor.Task{
		ID:    "task-1",
		Type:  supervisor.TaskTypeMarketTrend,
		Query: "analyze smartphone market",
	}

	_, err := agent.Execute(context.Background(), task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "market analysis failed")
}

func TestCompetitorAgent_Execute(t *testing.T) {
	mockLLM := &llm.MockLLM{
		RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
			return &llm.ChatResponse{Content: "Competitor analysis complete"}, nil
		},
	}

	agent := NewCompetitorAgent(mockLLM)
	task := &supervisor.Task{
		ID:    "task-2",
		Type:  supervisor.TaskTypeCompetitor,
		Query: "compare laptop competitors",
	}

	result, err := agent.Execute(context.Background(), task)

	require.NoError(t, err)
	compResult, ok := result.(*supervisor.CompetitorResult)
	require.True(t, ok)
	assert.NotEmpty(t, compResult.Analysis)
	assert.NotEmpty(t, compResult.Competitors)
}

func TestSupplyChainAgent_Execute(t *testing.T) {
	mockLLM := &llm.MockLLM{
		RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
			return &llm.ChatResponse{Content: "Supply chain stable"}, nil
		},
	}

	agent := NewSupplyChainAgent(mockLLM)
	task := &supervisor.Task{
		ID:    "task-3",
		Type:  supervisor.TaskTypeSupplyChain,
		Query: "analyze supply chain risks",
	}

	result, err := agent.Execute(context.Background(), task)

	require.NoError(t, err)
	scResult, ok := result.(*supervisor.SupplyChainResult)
	require.True(t, ok)
	assert.Equal(t, "Medium", scResult.RiskLevel)
	assert.NotEmpty(t, scResult.Factors)
}

func TestAgents_Name(t *testing.T) {
	marketAgent := NewMarketAgent(nil)
	assert.Equal(t, "MarketAgent", marketAgent.Name())

	competitorAgent := NewCompetitorAgent(nil)
	assert.Equal(t, "CompetitorAgent", competitorAgent.Name())

	supplyChainAgent := NewSupplyChainAgent(nil)
	assert.Equal(t, "SupplyChainAgent", supplyChainAgent.Name())
}
