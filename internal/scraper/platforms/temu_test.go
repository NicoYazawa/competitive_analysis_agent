package platforms

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockTemuScraper implements PlatformScraper for Temu tests
type mockTemuScraper struct {
	platform   Platform
	searchURL  string
	productURL string
	searchData []*CompetitorData
	product    *CompetitorData
	searchErr  error
	productErr error
}

func (m *mockTemuScraper) Platform() Platform                      { return m.platform }
func (m *mockTemuScraper) BuildSearchURL(query string) string       { return m.searchURL }
func (m *mockTemuScraper) BuildProductURL(productID string) string { return m.productURL }
func (m *mockTemuScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	return m.searchData, m.searchErr
}
func (m *mockTemuScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	return m.searchData, m.searchErr
}
func (m *mockTemuScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
	return m.product, m.productErr
}
func (m *mockTemuScraper) Close() error { return nil }

// Ensure mockTemuScraper implements PlatformScraper
var _ PlatformScraper = (*mockTemuScraper)(nil)

func TestTemuScraper_BuildSearchURL(t *testing.T) {
	tests := []struct {
		query string
		want  string
	}{
		{"test", "https://www.temu.com/search?q=test"},
		{"test query", "https://www.temu.com/search?q=test+query"},
		{"  spaces  ", "https://www.temu.com/search?q=++spaces++"},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			s := &TemuScraper{}
			got := s.BuildSearchURL(tt.query)
			if got != tt.want {
				t.Errorf("BuildSearchURL(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}

func TestTemuScraper_BuildProductURL(t *testing.T) {
	tests := []struct {
		productID string
		want     string
	}{
		{"123456789", "https://www.temu.com/item/123456789.html"},
		{"abcdefgh", "https://www.temu.com/item/abcdefgh.html"},
	}

	for _, tt := range tests {
		t.Run(tt.productID, func(t *testing.T) {
			s := &TemuScraper{}
			got := s.BuildProductURL(tt.productID)
			if got != tt.want {
				t.Errorf("BuildProductURL(%q) = %v, want %v", tt.productID, got, tt.want)
			}
		})
	}
}

func TestTemuScraper_Platform(t *testing.T) {
	s := &TemuScraper{}
	if s.Platform() != PlatformTemu {
		t.Errorf("Platform() = %v, want %v", s.Platform(), PlatformTemu)
	}
}

func TestTemuScraper_Interface(t *testing.T) {
	var s PlatformScraper = &TemuScraper{}
	if s.Platform() != PlatformTemu {
		t.Errorf("Platform() = %v, want %v", s.Platform(), PlatformTemu)
	}
	if s.BuildSearchURL("test") == "" {
		t.Error("BuildSearchURL returned empty string")
	}
	if s.BuildProductURL("123") == "" {
		t.Error("BuildProductURL returned empty string")
	}
}

func TestTemuScraper_MockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Temu Test</title></head>
<body>
<div class="search-result-list">
  <div class="search-result-item">
    <a href="/item/123456789.html"><span class="goods-title">Test Product</span></a>
    <span class="goods-price">$19.99</span>
    <span class="seller-name">Test Seller</span>
  </div>
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

func TestTemuScraper_URLBuilding(t *testing.T) {
	s := &TemuScraper{}

	// Test query with special characters
	url := s.BuildSearchURL("fashion")
	expected := "https://www.temu.com/search?q=fashion"
	if url != expected {
		t.Errorf("BuildSearchURL() = %v, want %v", url, expected)
	}

	// Test product URL
	productURL := s.BuildProductURL("987654321")
	expectedProduct := "https://www.temu.com/item/987654321.html"
	if productURL != expectedProduct {
		t.Errorf("BuildProductURL() = %v, want %v", productURL, expectedProduct)
	}
}
