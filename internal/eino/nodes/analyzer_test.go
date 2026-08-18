package nodes

import (
	"context"
	"errors"
	"testing"

	"competitive-analysis-agent/internal/llm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzerNode_WithLLM(t *testing.T) {
	mockLLM := &llm.MockLLM{
		RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
			return &llm.ChatResponse{Content: "LLM analysis result"}, nil
		},
	}

	node := NewAnalyzerNode(mockLLM)

	input := &AnalyzerInput{
		Query:        "test query",
		RawData:      "raw data content",
		AnalysisType: "market",
	}

	output, err := node.Invoke(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, "LLM analysis result", output.Analysis)
}

func TestAnalyzerNode_LLMError(t *testing.T) {
	mockLLM := &llm.MockLLM{
		ErrFn: func(ctx context.Context, req *llm.ChatRequest) error {
			return errors.New("LLM error")
		},
	}

	node := NewAnalyzerNode(mockLLM)

	input := &AnalyzerInput{
		Query:   "test",
		RawData: "data",
	}

	_, err := node.Invoke(context.Background(), input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "analysis failed")
}

func TestAnalyzerNode_EmptyInput(t *testing.T) {
	mockLLM := &llm.MockLLM{
		RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
			return &llm.ChatResponse{Content: "Mock analysis for empty input"}, nil
		},
	}
	node := NewAnalyzerNode(mockLLM)

	output, err := node.Invoke(context.Background(), &AnalyzerInput{})

	require.NoError(t, err)
	assert.Contains(t, output.Analysis, "Mock analysis for empty input")
	assert.NotEmpty(t, output.Insights)
}
