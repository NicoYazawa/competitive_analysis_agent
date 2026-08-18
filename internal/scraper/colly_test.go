package scraper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCollyScraper(t *testing.T) {
	scraper := NewCollyScraper()
	assert.NotNil(t, scraper)
}

func TestCollyScraper_SetProxies(t *testing.T) {
	scraper := NewCollyScraper()
	proxies := []string{"proxy1:8080", "proxy2:8080"}
	scraper.SetProxies(proxies)
	assert.Equal(t, proxies, scraper.proxies)
}

func TestCollyScraper_getNextProxy(t *testing.T) {
	scraper := NewCollyScraper()
	proxies := []string{"proxy1:8080", "proxy2:8080", "proxy3:8080"}
	scraper.SetProxies(proxies)

	// 测试轮换
	assert.Equal(t, "proxy1:8080", scraper.getNextProxy())
	assert.Equal(t, "proxy2:8080", scraper.getNextProxy())
	assert.Equal(t, "proxy3:8080", scraper.getNextProxy())
	assert.Equal(t, "proxy1:8080", scraper.getNextProxy()) // 循环
}

func TestCollyScraper_ScrapePage(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Test Product</title></head>
		<body>
			<span class="a-price-whole">99</span>
			<span class="a-icon-alt">4.5 out of 5</span>
			<span class="a-size-base">1,234 reviews</span>
			<div id="productDescription">Test description content</div>
		</body>
		</html>
		`
		w.Write([]byte(html))
	}))
	defer server.Close()

	scraper := NewCollyScraper()
	result, err := scraper.ScrapePage(context.Background(), server.URL)

	require.NoError(t, err)
	assert.Equal(t, "Test Product", result.Title)
	assert.Contains(t, result.Price, "99")
	assert.Contains(t, result.Rating, "4.5")
	assert.Contains(t, result.Reviews, "1,234")
	assert.Contains(t, result.Content, "Test description")
}

func TestCollyScraper_ScrapePage_Error(t *testing.T) {
	scraper := NewCollyScraper()
	_, err := scraper.ScrapePage(context.Background(), "http://invalid-host-that-does-not-exist.local")
	assert.Error(t, err)
}

func TestCollyScraper_ScrapeMultiple(t *testing.T) {
	// 创建两个测试服务器
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><title>Product 1</title><span class="a-price-whole">10</span></html>`))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><title>Product 2</title><span class="a-price-whole">20</span></html>`))
	}))
	defer server2.Close()

	scraper := NewCollyScraper()
	results := scraper.ScrapeMultiple(context.Background(), []string{server1.URL, server2.URL})

	assert.Len(t, results, 2)
	assert.Contains(t, results[0].Title, "Product 1")
	assert.Contains(t, results[1].Title, "Product 2")
}

func TestParsePrice(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{"simple dollar", "$99", 99.0, false},
		{"with cents", "$99.99", 99.99, false},
		{"with comma", "$1,234.56", 1234.56, false},
		{"no dollar", "99", 99.0, false},
		{"empty string", "", 0, true},
		{"invalid", "abc", 0, true},
		{"only spaces", "   ", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePrice(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseRating(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{"simple", "4.5", 4.5, false},
		{"out of 5", "4.5 out of 5", 4.5, false},
		{"with text", "Rated 4.5 stars", 4.5, false},
		{"integer", "5", 5.0, false},
		{"empty", "", 0, true},
		{"invalid", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRating(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseReviewCount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		hasError bool
	}{
		{"simple number", "1234", 1234, false},
		{"with comma", "1,234", 1234, false},
		{"with comma 2", "12,345", 12345, false},
		{"text with numbers", "1,234 reviews", 1234, false},
		{"zero", "0", 0, false},
		{"empty", "", 0, true},
		{"only letters", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseReviewCount(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestScrapeResult_Content(t *testing.T) {
	result := &ScrapResult{
		URL:     "https://example.com",
		Title:   "Test",
		Content: "Content",
		RawData: map[string]string{"key": "value"},
	}

	assert.Equal(t, "https://example.com", result.URL)
	assert.Equal(t, "Test", result.Title)
	assert.Equal(t, "Content", result.Content)
	assert.Equal(t, "value", result.RawData["key"])
}

func TestCollyScraper_Close(t *testing.T) {
	scraper := NewCollyScraper()
	assert.NotPanics(t, func() {
		scraper.Close()
	})
}
