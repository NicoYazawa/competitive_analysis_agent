package eino

import (
	"context"
	"testing"

	"competitive-analysis-agent/internal/llm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraph_Execute(t *testing.T) {
	mockLLM := &llm.MockLLM{
		RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
			return &llm.ChatResponse{Content: "Mock analysis result"}, nil
		},
	}
	graph := NewGraph(mockLLM)

	result, err := graph.Execute(context.Background(), "test query")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Strategy)
	assert.NotEmpty(t, result.Recommendations)
	assert.NotEmpty(t, result.AnalysisResults)
}

func TestGraph_Execute_DAGFlow(t *testing.T) {
	mockLLM := &llm.MockLLM{
		RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
			return &llm.ChatResponse{Content: "Mock analysis for: " + req.Messages[0].Content}, nil
		},
	}
	graph := NewGraph(mockLLM)

	result, err := graph.Execute(context.Background(), "smartphone pricing strategy")

	require.NoError(t, err)
	// Verify all stages executed
	assert.Contains(t, result.AnalysisResults[0], "smartphone pricing strategy")
	assert.Contains(t, result.Strategy, "smartphone pricing strategy")
	assert.Len(t, result.Recommendations, 3)
}
