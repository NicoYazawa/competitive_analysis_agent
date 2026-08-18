package platforms

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockAliExpressScraper implements PlatformScraper for AliExpress tests
type mockAliExpressScraper struct {
	platform   Platform
	searchURL  string
	productURL string
	searchData []*CompetitorData
	product    *CompetitorData
	searchErr  error
	productErr error
}

func (m *mockAliExpressScraper) Platform() Platform                      { return m.platform }
func (m *mockAliExpressScraper) BuildSearchURL(query string) string       { return m.searchURL }
func (m *mockAliExpressScraper) BuildProductURL(productID string) string { return m.productURL }
func (m *mockAliExpressScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	return m.searchData, m.searchErr
}
func (m *mockAliExpressScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	return m.searchData, m.searchErr
}
func (m *mockAliExpressScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
	return m.product, m.productErr
}
func (m *mockAliExpressScraper) Close() error { return nil }

// Ensure mockAliExpressScraper implements PlatformScraper
var _ PlatformScraper = (*mockAliExpressScraper)(nil)

func TestAliExpressScraper_BuildSearchURL(t *testing.T) {
	tests := []struct {
		query string
		want  string
	}{
		{"test", "https://www.aliexpress.com/wholesale?SearchText=test"},
		{"test query", "https://www.aliexpress.com/wholesale?SearchText=test+query"},
		{"  spaces  ", "https://www.aliexpress.com/wholesale?SearchText=++spaces++"},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			s := &AliExpressScraper{}
			got := s.BuildSearchURL(tt.query)
			if got != tt.want {
				t.Errorf("BuildSearchURL(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}

func TestAliExpressScraper_BuildProductURL(t *testing.T) {
	tests := []struct {
		productID string
		want     string
	}{
		{"123456789", "https://www.aliexpress.com/item/123456789.html"},
		{"abcdefgh", "https://www.aliexpress.com/item/abcdefgh.html"},
	}

	for _, tt := range tests {
		t.Run(tt.productID, func(t *testing.T) {
			s := &AliExpressScraper{}
			got := s.BuildProductURL(tt.productID)
			if got != tt.want {
				t.Errorf("BuildProductURL(%q) = %v, want %v", tt.productID, got, tt.want)
			}
		})
	}
}

func TestAliExpressScraper_Platform(t *testing.T) {
	s := &AliExpressScraper{}
	if s.Platform() != PlatformAliExpress {
		t.Errorf("Platform() = %v, want %v", s.Platform(), PlatformAliExpress)
	}
}

func TestAliExpressScraper_Interface(t *testing.T) {
	var s PlatformScraper = &AliExpressScraper{}
	if s.Platform() != PlatformAliExpress {
		t.Errorf("Platform() = %v, want %v", s.Platform(), PlatformAliExpress)
	}
	if s.BuildSearchURL("test") == "" {
		t.Error("BuildSearchURL returned empty string")
	}
	if s.BuildProductURL("123") == "" {
		t.Error("BuildProductURL returned empty string")
	}
}

func TestAliExpressScraper_MockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>AliExpress Test</title></head>
<body>
<div class="list-item">
  <a href="/item/123456789.html"><span class="item-title">Test Product</span></a>
  <span class="price">$19.99</span>
  <span class="orders-count">100 orders</span>
</div>
</body>
</html>`))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("HTTP GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestAliExpressScraper_URLBuilding(t *testing.T) {
	s := &AliExpressScraper{}

	// Test query with special characters
	url := s.BuildSearchURL("electronics")
	expected := "https://www.aliexpress.com/wholesale?SearchText=electronics"
	if url != expected {
		t.Errorf("BuildSearchURL() = %v, want %v", url, expected)
	}

	// Test product URL
	productURL := s.BuildProductURL("100001")
	expectedProduct := "https://www.aliexpress.com/item/100001.html"
	if productURL != expectedProduct {
		t.Errorf("BuildProductURL() = %v, want %v", productURL, expectedProduct)
	}
}
