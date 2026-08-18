package agents

import (
	"context"
	"fmt"

	"competitive-analysis-agent/internal/llm"
	"competitive-analysis-agent/internal/supervisor"
)

// SupplyChainAgent analyzes supply chain data.
type SupplyChainAgent struct {
	llmClient llm.ChatModel
}

// NewSupplyChainAgent creates a new SupplyChainAgent.
func NewSupplyChainAgent(llmClient llm.ChatModel) *SupplyChainAgent {
	return &SupplyChainAgent{llmClient: llmClient}
}

// Name returns the agent name.
func (a *SupplyChainAgent) Name() string {
	return "SupplyChainAgent"
}

// Execute performs supply chain analysis.
func (a *SupplyChainAgent) Execute(ctx context.Context, task *supervisor.Task) (any, error) {
	prompt := fmt.Sprintf("Analyze supply chain for: %s", task.Query)

	resp, err := a.llmClient.Chat(ctx, &llm.ChatRequest{
		Model: "qwen-max",
		Messages: []llm.ChatMessage{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("supply chain analysis failed: %w", err)
	}

	return &supervisor.SupplyChainResult{
		Status:    resp.Content,
		RiskLevel: "Medium",
		Factors:   []string{"Shipping cost", "Supplier reliability", "Inventory level"},
	}, nil
}
