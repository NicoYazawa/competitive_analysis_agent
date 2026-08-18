package llm

import (
	"fmt"

	"competitive-analysis-agent/internal/config"
)

// ProviderFactory creates LLM providers based on configuration.
type ProviderFactory struct{}

// NewProviderFactory creates a new ProviderFactory.
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{}
}

// CreateProvider creates a ChatModel provider based on the provider type.
func (f *ProviderFactory) CreateProvider(cfg *config.LLMConfig) (ChatModel, error) {
	switch config.LLMProviderType(cfg.Provider) {
	case config.LLMProviderQwen:
		model := cfg.Models["qwen"]
		if model == "" {
			model = "qwen-max"
		}
		provider := NewOpenAICompatProvider(cfg.APIKey, cfg.BaseURL, model, "qwen")
		// 百炼 uses /compatible-mode/v1/chat/completions
		if cfg.BaseURL == "https://dashscope.aliyuncs.com" {
			provider.SetAPIPath("/compatible-mode/v1/chat/completions")
		}
		return provider, nil

	case config.LLMProviderDeepSeek:
		model := cfg.Models["deepseek"]
		if model == "" {
			model = "deepseek-chat"
		}
		return NewOpenAICompatProvider(cfg.APIKey, cfg.BaseURL, model, "deepseek"), nil

	case config.LLMProviderMiniMax:
		model := cfg.Models["minimax"]
		if model == "" {
			model = "abab6-chat"
		}
		provider := NewOpenAICompatProvider(cfg.APIKey, cfg.BaseURL, model, "minimax")
		// MiniMax might need group_id header
		return provider, nil

	case config.LLMProviderMimiMax:
		model := cfg.Models["mimimax"]
		if model == "" {
			model = "mimimax-chat"
		}
		return NewOpenAICompatProvider(cfg.APIKey, cfg.BaseURL, model, "mimimax"), nil

	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", cfg.Provider)
	}
}

// CreateProviderByName creates a provider by name with explicit parameters.
func (f *ProviderFactory) CreateProviderByName(providerType, apiKey, baseURL, model string) (ChatModel, error) {
	switch config.LLMProviderType(providerType) {
	case config.LLMProviderQwen:
		return NewOpenAICompatProvider(apiKey, baseURL, model, "qwen"), nil
	case config.LLMProviderDeepSeek:
		return NewOpenAICompatProvider(apiKey, baseURL, model, "deepseek"), nil
	case config.LLMProviderMiniMax:
		return NewOpenAICompatProvider(apiKey, baseURL, model, "minimax"), nil
	case config.LLMProviderMimiMax:
		return NewOpenAICompatProvider(apiKey, baseURL, model, "mimimax"), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", providerType)
	}
}
