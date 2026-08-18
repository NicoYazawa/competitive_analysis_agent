package scraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractJSONFromPage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		checkFn     func(map[string]interface{})
	}{
		{
			name:        "valid json",
			input:       `{"key": "value", "number": 123}`,
			expectError: false,
			checkFn: func(m map[string]interface{}) {
				assert.Equal(t, "value", m["key"])
				assert.Equal(t, 123.0, m["number"])
			},
		},
		{
			name:        "nested json",
			input:       `{"outer": {"inner": "value"}}`,
			expectError: false,
			checkFn: func(m map[string]interface{}) {
				outer, ok := m["outer"].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "value", outer["inner"])
			},
		},
		{
			name:        "invalid json",
			input:       `{invalid}`,
			expectError: true,
		},
		{
			name:        "empty json",
			input:       `{}`,
			expectError: false,
			checkFn: func(m map[string]interface{}) {
				assert.Empty(t, m)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractJSONFromPage(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkFn != nil {
					tt.checkFn(result)
				}
			}
		})
	}
}

func TestCleanText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal", "Hello World", "Hello World"},
		{"multiple spaces", "Hello   World", "Hello World"},
		{"tabs and newlines", "Hello\t\nWorld", "Hello World"},
		{"leading trailing spaces", "  Hello  ", "Hello"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanText(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRodScrapResult(t *testing.T) {
	result := &RodScrapResult{
		URL:     "https://example.com",
		Title:   "Test Product",
		Price:   "99.99",
		Rating:  "4.5",
		Reviews: "1234",
		Seller:  "Test Seller",
		Content: "Test content",
		RawData: map[string]string{"key": "value"},
		Error:   nil,
		RawHTML: "<html></html>",
	}

	assert.Equal(t, "https://example.com", result.URL)
	assert.Equal(t, "Test Product", result.Title)
	assert.Equal(t, "99.99", result.Price)
	assert.Equal(t, "4.5", result.Rating)
	assert.Equal(t, "1234", result.Reviews)
	assert.Equal(t, "Test Seller", result.Seller)
	assert.Equal(t, "Test content", result.Content)
	assert.Equal(t, "value", result.RawData["key"])
	assert.Nil(t, result.Error)
	assert.Equal(t, "<html></html>", result.RawHTML)
}

func TestRodScrapResult_WithError(t *testing.T) {
	result := &RodScrapResult{
		URL:   "https://example.com",
		Title: "",
		Error: assert.AnError,
	}

	assert.Equal(t, "https://example.com", result.URL)
	assert.Error(t, result.Error)
}

func TestDefaultTimeout(t *testing.T) {
	// 测试默认超时常量
	assert.Equal(t, 30_000_000_000, defaultTimeout)
}
