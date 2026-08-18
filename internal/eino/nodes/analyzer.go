package nodes

import (
	"context"
	"errors"
	"fmt"

	"competitive-analysis-agent/internal/llm"
)

// ErrAnalyzerNilClient is returned when AnalyzerNode is invoked without an LLM client.
var ErrAnalyzerNilClient = errors.New("analyzer node requires an LLM client, got nil")

// AnalyzerInput represents input to the analyzer node.
type AnalyzerInput struct {
	Query        string `json:"query"`
	RawData      string `json:"raw_data"`
	AnalysisType string `json:"analysis_type"`
}

// AnalyzerOutput represents output from the analyzer node.
type AnalyzerOutput struct {
	Analysis string   `json:"analysis"`
	Insights []string `json:"insights"`
}

// AnalyzerNode analyzes data using LLM.
type AnalyzerNode struct {
	llmClient llm.ChatModel
}

// NewAnalyzerNode creates a new AnalyzerNode. Pass a real or mock LLM client.
func NewAnalyzerNode(llmClient llm.ChatModel) *AnalyzerNode {
	return &AnalyzerNode{llmClient: llmClient}
}

// Invoke executes the analyzer node.
func (n *AnalyzerNode) Invoke(ctx context.Context, input *AnalyzerInput) (*AnalyzerOutput, error) {
	if n.llmClient == nil {
		return nil, ErrAnalyzerNilClient
	}

	prompt := fmt.Sprintf("Analyze the following data for: %s\n\nData: %s", input.Query, input.RawData)

	resp, err := n.llmClient.Chat(ctx, &llm.ChatRequest{
		Model: "qwen-max",
		Messages: []llm.ChatMessage{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}

	return &AnalyzerOutput{
		Analysis: resp.Content,
		Insights: []string{"Auto-generated insight"},
	}, nil
}
