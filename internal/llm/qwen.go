package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OpenAICompatProvider implements ChatModel using OpenAI-compatible API.
type OpenAICompatProvider struct {
	APIKey    string
	BaseURL   string
	Model     string
	Provider  string // provider name for logging
	APIPath   string // custom API path, default "/v1/chat/completions"
	client    *http.Client
	headers   map[string]string // additional headers
}

// NewOpenAICompatProvider creates a new OpenAI-compatible Provider.
func NewOpenAICompatProvider(apiKey, baseURL, model, provider string) *OpenAICompatProvider {
	return &OpenAICompatProvider{
		APIKey:   apiKey,
		BaseURL:  baseURL,
		Model:    model,
		Provider: provider,
		APIPath:  "/v1/chat/completions",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}
}

// NewQwenProvider creates a new QwenProvider (backward compatible).
func NewQwenProvider(apiKey, baseURL, model string) *OpenAICompatProvider {
	return NewOpenAICompatProvider(apiKey, baseURL, model, "qwen")
}

// SetHeader sets a custom header for the requests.
func (p *OpenAICompatProvider) SetHeader(key, value string) {
	p.headers[key] = value
}

// SetAPIPath sets a custom API path (e.g., "/api/v1/chat/completions" for 百炼).
func (p *OpenAICompatProvider) SetAPIPath(path string) {
	p.APIPath = path
}

// Chat implements ChatModel interface.
func (p *OpenAICompatProvider) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	url := fmt.Sprintf("%s%s", p.BaseURL, p.APIPath)

	httpReq := struct {
		Model    string        `json:"model"`
		Messages []ChatMessage `json:"messages"`
	}{
		Model:    p.getModel(req),
		Messages: req.Messages,
	}

	body, err := json.Marshal(httpReq)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to marshal request: %w", p.Provider, err)
	}

	httpReq2, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to create request: %w", p.Provider, err)
	}

	httpReq2.Header.Set("Content-Type", "application/json")
	httpReq2.Header.Set("Authorization", "Bearer "+p.APIKey)

	for k, v := range p.headers {
		httpReq2.Header.Set(k, v)
	}

	resp, err := p.client.Do(httpReq2)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to send request: %w", p.Provider, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[%s] unexpected status code: %d", p.Provider, resp.StatusCode)
	}

	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("[%s] failed to decode response: %w", p.Provider, err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("[%s] no choices in response", p.Provider)
	}

	return &ChatResponse{
		Content: chatResp.Choices[0].Message.Content,
	}, nil
}

// getModel returns the model to use, defaulting to provider's model if not specified in request.
func (p *OpenAICompatProvider) getModel(req *ChatRequest) string {
	if req.Model != "" {
		return req.Model
	}
	return p.Model
}
