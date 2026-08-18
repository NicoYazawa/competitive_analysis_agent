package agents

import (
	"context"
	"fmt"

	"competitive-analysis-agent/internal/llm"
	"competitive-analysis-agent/internal/supervisor"
)

// CompetitorAgent analyzes competitor data.
type CompetitorAgent struct {
	llmClient llm.ChatModel
}

// NewCompetitorAgent creates a new CompetitorAgent.
func NewCompetitorAgent(llmClient llm.ChatModel) *CompetitorAgent {
	return &CompetitorAgent{llmClient: llmClient}
}

// Name returns the agent name.
func (a *CompetitorAgent) Name() string {
	return "CompetitorAgent"
}

// Execute performs competitor analysis.
func (a *CompetitorAgent) Execute(ctx context.Context, task *supervisor.Task) (any, error) {
	prompt := fmt.Sprintf("Analyze competitors for: %s", task.Query)

	resp, err := a.llmClient.Chat(ctx, &llm.ChatRequest{
		Model: "qwen-max",
		Messages: []llm.ChatMessage{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("competitor analysis failed: %w", err)
	}

	return &supervisor.CompetitorResult{
		Analysis: resp.Content,
		Competitors: []supervisor.CompetitorInsight{
			{Name: "Competitor A", Strength: "Strong brand", Weakness: "High price", Strategy: "Premium positioning"},
		},
	}, nil
}
