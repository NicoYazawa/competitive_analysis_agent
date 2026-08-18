package llm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockLLM_CalledWithRecording(t *testing.T) {
	mock := &MockLLM{
		RespFn: func(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
			return &ChatResponse{Content: "response: " + req.Messages[0].Content}, nil
		},
	}

	// Make multiple calls
	mock.Chat(context.Background(), &ChatRequest{Model: "model1", Messages: []ChatMessage{{Role: "user", Content: "msg1"}}})
	mock.Chat(context.Background(), &ChatRequest{Model: "model2", Messages: []ChatMessage{{Role: "user", Content: "msg2"}}})

	assert.Len(t, mock.CalledWith, 2)
	assert.Equal(t, "model1", mock.CalledWith[0].Model)
	assert.Equal(t, "model2", mock.CalledWith[1].Model)
	assert.Equal(t, "msg1", mock.CalledWith[0].Messages[0].Content)
	assert.Equal(t, "msg2", mock.CalledWith[1].Messages[0].Content)
}

func TestMockLLM_EmptyMessages(t *testing.T) {
	mock := &MockLLM{}

	resp, err := mock.Chat(context.Background(), &ChatRequest{
		Model:    "test",
		Messages: []ChatMessage{},
	})

	require.NoError(t, err)
	assert.Equal(t, "mock response", resp.Content)
}

func TestMockLLM_ErrFnReturningNilDoesNotBlockRespFn(t *testing.T) {
	mock := &MockLLM{
		ErrFn: func(ctx context.Context, req *ChatRequest) error {
			return nil // no error
		},
		RespFn: func(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
			return &ChatResponse{Content: "response via RespFn"}, nil
		},
	}

	resp, err := mock.Chat(context.Background(), &ChatRequest{})

	require.NoError(t, err)
	assert.Equal(t, "response via RespFn", resp.Content)
}
