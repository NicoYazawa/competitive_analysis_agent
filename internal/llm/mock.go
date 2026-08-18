package llm

import (
	"context"
)

// MockLLM is a mock implementation of ChatModel for testing.
type MockLLM struct {
	RespFn     func(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	ErrFn      func(ctx context.Context, req *ChatRequest) error
	CalledWith []*ChatRequest
}

// Chat implements ChatModel interface.
func (m *MockLLM) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	m.CalledWith = append(m.CalledWith, req)
	if m.ErrFn != nil {
		if err := m.ErrFn(ctx, req); err != nil {
			return nil, err
		}
	}
	if m.RespFn != nil {
		return m.RespFn(ctx, req)
	}
	return &ChatResponse{Content: "mock response"}, nil
}
