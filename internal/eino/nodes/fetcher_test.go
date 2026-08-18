package nodes

import (
	"context"
	"testing"

	"competitive-analysis-agent/internal/llm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetcherNode_Invoke(t *testing.T) {
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

func TestAnalyzerNode_Invoke(t *testing.T) {
	mockLLM := &llm.MockLLM{
		RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
			return &llm.ChatResponse{Content: "Mock analysis for: " + req.Messages[0].Content}, nil
		},
	}
	node := NewAnalyzerNode(mockLLM)

	input := &AnalyzerInput{
		Query:   "test query",
		RawData: "raw data for analysis",
	}

	output, err := node.Invoke(context.Background(), input)

	require.NoError(t, err)
	assert.Contains(t, output.Analysis, "Mock analysis")
	assert.NotEmpty(t, output.Insights)
}

func TestSynthesizerNode_Invoke(t *testing.T) {
	node := NewSynthesizerNode()

	input := &SynthesizerInput{
		Query: "smartphone strategy",
		AnalysisResults: []string{
			"Analysis 1: Price trends",
			"Analysis 2: Competitor landscape",
		},
	}

	output, err := node.Invoke(context.Background(), input)

	require.NoError(t, err)
	assert.Contains(t, output.FinalStrategy, "smartphone strategy")
	assert.Len(t, output.Recommendations, 3)
}

func TestSynthesizerNode_Invoke_EmptyResults(t *testing.T) {
	node := NewSynthesizerNode()

	input := &SynthesizerInput{
		Query:           "test",
		AnalysisResults: []string{},
	}

	output, err := node.Invoke(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, "No analysis results to synthesize", output.FinalStrategy)
}
