package llm

import (
	"testing"

	"competitive-analysis-agent/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderFactory_CreateProvider_Qwen(t *testing.T) {
	factory := NewProviderFactory()

	cfg := &config.LLMConfig{
		Provider: "qwen",
		APIKey:   "test-key",
		BaseURL:  "https://api.qwen.com",
		Models:   map[string]string{"qwen": "qwen-max"},
	}

	provider, err := factory.CreateProvider(cfg)

	require.NoError(t, err)
	assert.NotNil(t, provider)
}

func TestProviderFactory_CreateProvider_DeepSeek(t *testing.T) {
	factory := NewProviderFactory()

	cfg := &config.LLMConfig{
		Provider: "deepseek",
		APIKey:   "test-key",
		BaseURL:  "https://api.deepseek.com",
		Models:   map[string]string{"deepseek": "deepseek-chat"},
	}

	provider, err := factory.CreateProvider(cfg)

	require.NoError(t, err)
	assert.NotNil(t, provider)
}

func TestProviderFactory_CreateProvider_MiniMax(t *testing.T) {
	factory := NewProviderFactory()

	cfg := &config.LLMConfig{
		Provider: "minimax",
		APIKey:   "test-key",
		BaseURL:  "https://api.minimax.chat",
		Models:   map[string]string{"minimax": "abab6-chat"},
	}

	provider, err := factory.CreateProvider(cfg)

	require.NoError(t, err)
	assert.NotNil(t, provider)
}

func TestProviderFactory_CreateProvider_MimiMax(t *testing.T) {
	factory := NewProviderFactory()

	cfg := &config.LLMConfig{
		Provider: "mimimax",
		APIKey:   "test-key",
		BaseURL:  "https://api.mimimax.io",
		Models:   map[string]string{"mimimax": "mimimax-chat"},
	}

	provider, err := factory.CreateProvider(cfg)

	require.NoError(t, err)
	assert.NotNil(t, provider)
}

func TestProviderFactory_CreateProvider_DefaultModel(t *testing.T) {
	factory := NewProviderFactory()

	cfg := &config.LLMConfig{
		Provider: "qwen",
		APIKey:   "test-key",
		BaseURL:  "https://api.qwen.com",
		// No models specified, should use default
	}

	provider, err := factory.CreateProvider(cfg)

	require.NoError(t, err)
	assert.NotNil(t, provider)
}

func TestProviderFactory_CreateProvider_Unsupported(t *testing.T) {
	factory := NewProviderFactory()

	cfg := &config.LLMConfig{
		Provider: "unsupported",
		APIKey:   "test-key",
		BaseURL:  "https://api.example.com",
	}

	_, err := factory.CreateProvider(cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported LLM provider")
}

func TestProviderFactory_CreateProviderByName(t *testing.T) {
	factory := NewProviderFactory()

	tests := []struct {
		provider string
		model    string
	}{
		{"qwen", "qwen-max"},
		{"deepseek", "deepseek-chat"},
		{"minimax", "abab6-chat"},
		{"mimimax", "mimimax-chat"},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			provider, err := factory.CreateProviderByName(
				tt.provider,
				"test-key",
				"https://api.example.com",
				tt.model,
			)

			require.NoError(t, err)
			assert.NotNil(t, provider)
		})
	}
}

func TestProviderFactory_CreateProviderByName_Unsupported(t *testing.T) {
	factory := NewProviderFactory()

	_, err := factory.CreateProviderByName(
		"unknown",
		"test-key",
		"https://api.example.com",
		"some-model",
	)

	assert.Error(t, err)
}
