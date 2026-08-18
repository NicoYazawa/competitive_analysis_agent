package eino

import (
	"context"
	"testing"

	"competitive-analysis-agent/internal/llm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraph_Execute_DAGFlowWithMockLLM(t *testing.T) {
	mockLLM := &llm.MockLLM{
		RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
			return &llm.ChatResponse{Content: "analysis for: " + req.Messages[0].Content}, nil
		},
	}
	graph := NewGraph(mockLLM)

	result, err := graph.Execute(context.Background(), "smartphone pricing analysis")

	require.NoError(t, err)
	assert.NotNil(t, result)

	// Verify all nodes executed
	assert.Len(t, result.AnalysisResults, 1)
	assert.NotEmpty(t, result.AnalysisResults[0])

	// Verify synthesizer output
	assert.NotEmpty(t, result.Strategy)
	assert.Len(t, result.Recommendations, 3)

	// Verify recommendations are not empty
	for _, rec := range result.Recommendations {
		assert.NotEmpty(t, rec)
	}
}

func TestGraph_Execute_VerifyDataFlow(t *testing.T) {
	mockLLM := &llm.MockLLM{
		RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
			return &llm.ChatResponse{Content: "electronics market analysis"}, nil
		},
	}
	graph := NewGraph(mockLLM)

	result, err := graph.Execute(context.Background(), "electronics market")

	require.NoError(t, err)

	// Verify fetcher output flows to analyzer
	assert.NotEmpty(t, result.AnalysisResults)
	assert.Contains(t, result.AnalysisResults[0], "electronics market")

	// Verify analyzer output flows to synthesizer
	assert.Contains(t, result.Strategy, "electronics market")
}

func TestGraph_NewGraph(t *testing.T) {
	mockLLM := &llm.MockLLM{
		RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
			return &llm.ChatResponse{Content: "response"}, nil
		},
	}
	graph := NewGraph(mockLLM)

	assert.NotNil(t, graph.fetcher)
	assert.NotNil(t, graph.analyzer)
	assert.NotNil(t, graph.synthesizer)
}
