package llm

import (
	"context"
)

// ChatMessage represents a single message in a chat conversation.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a request to the LLM chat endpoint.
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatResponse represents a response from the LLM chat endpoint.
type ChatResponse struct {
	Content string `json:"content"`
}

// ChatModel is the interface for LLM chat models.
type ChatModel interface {
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
}
