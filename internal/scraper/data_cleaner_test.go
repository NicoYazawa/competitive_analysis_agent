package scraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDataCleaner(t *testing.T) {
	cleaner := NewDataCleaner()
	assert.NotNil(t, cleaner)
}

func TestDataCleaner_CleanText(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal text", "Hello World", "Hello World"},
		{"multiple spaces", "Hello   World", "Hello World"},
		{"leading spaces", "  Hello", "Hello"},
		{"trailing spaces", "Hello  ", "Hello"},
		{"tabs and newlines", "Hello\t\nWorld", "Hello World"},
		{"empty string", "", ""},
		{"chinese text", "你好世界", "你好世界"},
		{"mixed text", "Hello 世界 123", "Hello 世界 123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.CleanText(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataCleaner_CleanPlatform(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"amazon lowercase", "amazon", "amazon"},
		{"amazon with domain", "amazon.com", "amazon"},
		{"amazon www", "www.amazon.com", "amazon"},
		{"aliexpress", "aliexpress", "aliexpress"},
		{"aliexpress www", "www.aliexpress.com", "aliexpress"},
		{"ebay", "ebay", "ebay"},
		{"unknown", "unknown_platform", "unknown_platform"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.CleanPlatform(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataCleaner_CleanCurrency(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"USD", "USD", "USD"},
		{"US$", "US$", "USD"},
		{"dollar sign", "$", "USD"},
		{"EUR", "EUR", "EUR"},
		{"euro sign", "€", "EUR"},
		{"GBP", "GBP", "GBP"},
		{"pound sign", "£", "GBP"},
		{"CNY", "CNY", "CNY"},
		{"JPY", "JPY", "JPY"},
		{"unknown", "XYZ", "XYZ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.CleanCurrency(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataCleaner_CleanCompetitorData(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		name        string
		input       *CompetitorData
		expectError bool
		checkFields func(*CleanedCompetitorData)
	}{
		{
			name: "valid data",
			input: &CompetitorData{
				Name:              "Test Product",
				Platform:          "amazon",
				PlatformProductID: "B001",
				Price:             "$99.99",
				Currency:          "USD",
				Rating:            "4.5 out of 5",
				ReviewCount:       "1,234",
				SourceURL:         "https://amazon.com/dp/B001",
			},
			expectError: false,
			checkFields: func(c *CleanedCompetitorData) {
				assert.Equal(t, "Test Product", c.Name)
				assert.Equal(t, "amazon", c.Platform)
				assert.Equal(t, 99.99, c.Price)
				assert.Equal(t, 4.5, c.Rating)
				assert.Equal(t, 1234, c.ReviewCount)
			},
		},
		{
			name: "empty name",
			input: &CompetitorData{
				Name:     "",
				Platform: "amazon",
			},
			expectError: true,
		},
		{
			name: "empty platform",
			input: &CompetitorData{
				Name:     "Test",
				Platform: "",
			},
			expectError: true,
		},
		{
			name: "invalid price",
			input: &CompetitorData{
				Name:     "Test",
				Platform: "amazon",
				Price:    "invalid",
			},
			expectError: false, // 价格解析失败会设为0，但不返回错误
			checkFields: func(c *CleanedCompetitorData) {
				assert.Equal(t, 0.0, c.Price)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cleaner.CleanCompetitorData(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checkFields != nil {
					tt.checkFields(result)
				}
			}
		})
	}
}

func TestDataCleaner_RemoveHTMLTags(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple tags", "<p>Hello</p>", "Hello"},
		{"nested tags", "<div><p>Hello</p></div>", "Hello"},
		{"nbsp", "Hello&nbsp;World", "Hello World"},
		{"ampersand", "Tom & Jerry", "Tom & Jerry"},
		{"quotes", `He said "Hello"`, `He said "Hello"`},
		{"no tags", "Plain Text", "Plain Text"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.RemoveHTMLTags(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataCleaner_TruncateString(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"short string", "Hello", 10, "Hello"},
		{"exact length", "Hello", 5, "Hello"},
		{"long string", "Hello World", 8, "Hello..."},
		{"empty string", "", 5, ""},
		{"maxLen 0", "Hello", 0, ""},
		{"maxLen 1", "Hello", 1, "H"},
		{"maxLen 2", "Hello", 2, "He"},
		{"maxLen 3", "Hello", 3, "Hel"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.TruncateString(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataCleaner_ValidateURL(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		url      string
		expected bool
	}{
		{"https://amazon.com/dp/B001", true},
		{"http://example.com", true},
		{"https://example.com/path?query=1", true},
		{"invalid", false},
		{"ftp://example.com", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := cleaner.ValidateURL(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataCleaner_ExtractDomain(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		url      string
		expected string
	}{
		{"https://www.amazon.com/dp/B001", "www.amazon.com"},
		{"http://example.com/path", "example.com"},
		{"https://sub.example.com", "sub.example.com"},
		{"invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := cleaner.ExtractDomain(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataCleaner_NormalizeWhitespace(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single spaces", "Hello World", "Hello World"},
		{"multiple spaces", "Hello   World", "Hello World"},
		{"tabs", "Hello\tWorld", "Hello World"},
		{"newlines", "Hello\nWorld", "Hello World"},
		{"mixed", "Hello  \t\n  World", "Hello World"},
		{"leading trailing", "  Hello  ", "Hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.NormalizeWhitespace(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataCleaner_JSONToMap(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		name        string
		input       string
		expectedKey string
		expectedVal interface{}
		expectError bool
	}{
		{
			name:        "valid json",
			input:       `{"key": "value"}`,
			expectedKey: "key",
			expectedVal: "value",
			expectError: false,
		},
		{
			name:        "invalid json",
			input:       `{invalid}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cleaner.JSONToMap(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedVal, result[tt.expectedKey])
			}
		})
	}
}

func TestDataCleaner_MapToJSON(t *testing.T) {
	cleaner := NewDataCleaner()

	input := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	result, err := cleaner.MapToJSON(input)
	require.NoError(t, err)
	assert.Contains(t, result, `"key1"`)
	assert.Contains(t, result, `"value1"`)
	assert.Contains(t, result, `"key2"`)
}

func TestDataCleaner_SanitizeForDB(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal text", "normal text", "normal text"},
		{"with single quote", "Tom's", "Tom''s"},
		{"with backslash", "path\\file", "path\\\\file"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.SanitizeForDB(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataCleaner_ExtractProductID(t *testing.T) {
	cleaner := NewDataCleaner()

	tests := []struct {
		name     string
		platform string
		url      string
		expected string
	}{
		{"amazon asin", "amazon", "https://amazon.com/dp/B001XXXXX", "B001XXXXX"},
		{"aliexpress product", "aliexpress", "https://www.aliexpress.com/item/1234567890.html", "1234567890"},
		{"ebay item", "ebay", "https://www.ebay.com/itm/1234567890", "1234567890"},
		{"unknown platform", "unknown", "https://example.com/product/123", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleaner.ExtractProductID(tt.url, tt.platform)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataCleaner_CleanPrice(t *testing.T) {
	cleaner := NewDataCleaner()

	result, err := cleaner.CleanPrice("$99.99")
	require.NoError(t, err)
	assert.Equal(t, 99.99, result)
}

func TestDataCleaner_CleanRating(t *testing.T) {
	cleaner := NewDataCleaner()

	result, err := cleaner.CleanRating("4.5 out of 5")
	require.NoError(t, err)
	assert.Equal(t, 4.5, result)
}

func TestDataCleaner_CleanReviewCount(t *testing.T) {
	cleaner := NewDataCleaner()

	result, err := cleaner.CleanReviewCount("1,234 reviews")
	require.NoError(t, err)
	assert.Equal(t, 1234, result)
}
