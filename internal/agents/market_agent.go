package agents

import (
	"context"
	"fmt"

	"competitive-analysis-agent/internal/llm"
	"competitive-analysis-agent/internal/supervisor"
)

// MarketAgent analyzes market trends.
type MarketAgent struct {
	llmClient llm.ChatModel
}

// NewMarketAgent creates a new MarketAgent.
func NewMarketAgent(llmClient llm.ChatModel) *MarketAgent {
	return &MarketAgent{llmClient: llmClient}
}

// Name returns the agent name.
func (a *MarketAgent) Name() string {
	return "MarketAgent"
}

// Execute performs market trend analysis.
func (a *MarketAgent) Execute(ctx context.Context, task *supervisor.Task) (any, error) {
	prompt := fmt.Sprintf("Analyze market trends for: %s", task.Query)

	resp, err := a.llmClient.Chat(ctx, &llm.ChatRequest{
		Model: "qwen-max",
		Messages: []llm.ChatMessage{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("market analysis failed: %w", err)
	}

	return &supervisor.MarketTrendResult{
		Trend:         resp.Content,
		Opportunities: []string{"Opportunity 1", "Opportunity 2"},
		DemandSignal:  "Strong demand detected",
	}, nil
}
