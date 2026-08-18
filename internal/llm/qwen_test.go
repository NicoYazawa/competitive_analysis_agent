package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQwenProvider(t *testing.T) {
	provider := NewQwenProvider("test-key", "https://api.qwen.com", "qwen-max")

	assert.Equal(t, "test-key", provider.APIKey)
	assert.Equal(t, "https://api.qwen.com", provider.BaseURL)
	assert.Equal(t, "qwen-max", provider.Model)
	assert.NotNil(t, provider.client)
}

func TestQwenProvider_Chat_InvalidURL(t *testing.T) {
	provider := NewQwenProvider("test-key", "://invalid-url", "qwen-max")

	_, err := provider.Chat(context.Background(), &ChatRequest{
		Model:    "qwen-max",
		Messages: []ChatMessage{{Role: "user", Content: "hello"}},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create request")
}

func TestQwenProvider_Chat_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	provider := NewQwenProvider("test-key", server.URL, "qwen-max")

	_, err := provider.Chat(context.Background(), &ChatRequest{
		Model:    "qwen-max",
		Messages: []ChatMessage{{Role: "user", Content: "hello"}},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode response")
}

func TestQwenProvider_Chat_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"choices": []}`))
	}))
	defer server.Close()

	provider := NewQwenProvider("test-key", server.URL, "qwen-max")

	_, err := provider.Chat(context.Background(), &ChatRequest{
		Model:    "qwen-max",
		Messages: []ChatMessage{{Role: "user", Content: "hello"}},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no choices in response")
}

func TestQwenProvider_Chat_RequestTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	provider := NewQwenProvider("test-key", server.URL, "qwen-max")
	provider.client.Timeout = 50 * time.Millisecond

	_, err := provider.Chat(context.Background(), &ChatRequest{
		Model:    "qwen-max",
		Messages: []ChatMessage{{Role: "user", Content: "hello"}},
	})

	assert.Error(t, err)
}

func TestQwenProvider_Chat_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	provider := NewQwenProvider("test-key", server.URL, "qwen-max")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := provider.Chat(ctx, &ChatRequest{
		Model:    "qwen-max",
		Messages: []ChatMessage{{Role: "user", Content: "hello"}},
	})

	assert.Error(t, err)
}

func TestQwenProvider_Chat_VerifyRequestBody(t *testing.T) {
	var receivedBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"choices": [{"message": {"content": "test"}}]}`))
	}))
	defer server.Close()

	provider := NewQwenProvider("test-key", server.URL, "qwen-max")

	_, err := provider.Chat(context.Background(), &ChatRequest{
		Model: "custom-model",
		Messages: []ChatMessage{
			{Role: "system", Content: "You are helpful"},
			{Role: "user", Content: "Hello"},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "custom-model", receivedBody["model"])
	assert.Len(t, receivedBody["messages"], 2)
}
