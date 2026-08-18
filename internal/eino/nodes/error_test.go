package nodes

import (
	"context"
	"errors"
	"testing"

	"competitive-analysis-agent/internal/llm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetcherNode_Invoke_WithError(t *testing.T) {
	node := &FetcherNode{InjectError: ErrFetcherFailed}

	input := &FetcherInput{
		Query: "test",
		Type:  "market_data",
	}

	_, err := node.Invoke(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "fetcher failed")
}

func TestFetcherNode_Invoke_WithoutError(t *testing.T) {
	node := NewFetcherNode()

	input := &FetcherInput{
		Query: "smartphone trends",
		Type:  "market_data",
	}

	output, err := node.Invoke(context.Background(), input)

	require.NoError(t, err)
	assert.True(t, output.Found)
	assert.Contains(t, output.Data, "smartphone trends")
}

func TestSynthesizerNode_Invoke_WithError(t *testing.T) {
	node := &SynthesizerNode{InjectError: ErrSynthesisFailed}

	input := &SynthesizerInput{
		Query:           "test",
		AnalysisResults: []string{"result1"},
	}

	_, err := node.Invoke(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "synthesis failed")
}

func TestSynthesizerNode_Invoke_WithMultipleResults(t *testing.T) {
	node := NewSynthesizerNode()

	input := &SynthesizerInput{
		Query: "test query",
		AnalysisResults: []string{
			"Analysis 1: Price trends",
			"Analysis 2: Competitor landscape",
			"Analysis 3: Market demand",
		},
	}

	output, err := node.Invoke(context.Background(), input)

	require.NoError(t, err)
	assert.Contains(t, output.FinalStrategy, "3 analysis results")
	assert.Len(t, output.Recommendations, 3)
}

func TestErrFetcherFailed(t *testing.T) {
	assert.NotNil(t, ErrFetcherFailed)
	assert.Equal(t, "fetcher failed to fetch data", ErrFetcherFailed.Error())
}

func TestErrSynthesisFailed(t *testing.T) {
	assert.NotNil(t, ErrSynthesisFailed)
	assert.Equal(t, "synthesis failed", ErrSynthesisFailed.Error())
}

func TestAnalyzerNode_WithMockLLM(t *testing.T) {
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
	assert.NotEmpty(t, output.Insights)
}

func TestAnalyzerNode_WithMockLLMError(t *testing.T) {
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
