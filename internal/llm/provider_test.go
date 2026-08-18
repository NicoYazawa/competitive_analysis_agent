package llm

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockLLM_Chat(t *testing.T) {
	t.Run("successful response", func(t *testing.T) {
		mock := &MockLLM{
			RespFn: func(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
				return &ChatResponse{Content: "test response"}, nil
			},
		}

		resp, err := mock.Chat(context.Background(), &ChatRequest{
			Model: "test-model",
			Messages: []ChatMessage{
				{Role: "user", Content: "hello"},
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "test response", resp.Content)
		assert.Len(t, mock.CalledWith, 1)
		assert.Equal(t, "test-model", mock.CalledWith[0].Model)
	})

	t.Run("error response", func(t *testing.T) {
		mock := &MockLLM{
			ErrFn: func(ctx context.Context, req *ChatRequest) error {
				return errors.New("mock error")
			},
		}

		_, err := mock.Chat(context.Background(), &ChatRequest{})

		assert.Error(t, err)
		assert.Equal(t, "mock error", err.Error())
	})

	t.Run("default response when no functions set", func(t *testing.T) {
		mock := &MockLLM{}

		resp, err := mock.Chat(context.Background(), &ChatRequest{})

		require.NoError(t, err)
		assert.Equal(t, "mock response", resp.Content)
	})
}

func TestQwenProvider_Chat(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/v1/chat/completions")
			assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"choices": [{
					"message": {
						"content": "test response"
					}
				}]
			}`))
		}))
		defer server.Close()

		provider := NewQwenProvider("test-key", server.URL, "qwen-max")
		resp, err := provider.Chat(context.Background(), &ChatRequest{
			Model: "qwen-max",
			Messages: []ChatMessage{
				{Role: "user", Content: "test prompt"},
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "test response", resp.Content)
	})

	t.Run("error status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		provider := NewQwenProvider("test-key", server.URL, "qwen-max")
		_, err := provider.Chat(context.Background(), &ChatRequest{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "500")
	})
}
