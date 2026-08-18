package platforms

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockEbayScraper implements PlatformScraper for eBay tests
type mockEbayScraper struct {
	platform   Platform
	searchURL  string
	productURL string
	searchData []*CompetitorData
	product    *CompetitorData
	searchErr  error
	productErr error
}

func (m *mockEbayScraper) Platform() Platform                      { return m.platform }
func (m *mockEbayScraper) BuildSearchURL(query string) string       { return m.searchURL }
func (m *mockEbayScraper) BuildProductURL(productID string) string { return m.productURL }
func (m *mockEbayScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	return m.searchData, m.searchErr
}
func (m *mockEbayScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	return m.searchData, m.searchErr
}
func (m *mockEbayScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
	return m.product, m.productErr
}
func (m *mockEbayScraper) Close() error { return nil }

// Ensure mockEbayScraper implements PlatformScraper
var _ PlatformScraper = (*mockEbayScraper)(nil)

func TestEbayScraper_BuildSearchURL(t *testing.T) {
	tests := []struct {
		query string
		want  string
	}{
		{"test", "https://www.ebay.com/sch/i.html?_nkw=test"},
		{"test query", "https://www.ebay.com/sch/i.html?_nkw=test+query"},
		{"  spaces  ", "https://www.ebay.com/sch/i.html?_nkw=++spaces++"},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			s := &EbayScraper{}
			got := s.BuildSearchURL(tt.query)
			if got != tt.want {
				t.Errorf("BuildSearchURL(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}

func TestEbayScraper_BuildProductURL(t *testing.T) {
	tests := []struct {
		productID string
		want     string
	}{
		{"123456789", "https://www.ebay.com/itm/123456789"},
		{"abcdefgh", "https://www.ebay.com/itm/abcdefgh"},
	}

	for _, tt := range tests {
		t.Run(tt.productID, func(t *testing.T) {
			s := &EbayScraper{}
			got := s.BuildProductURL(tt.productID)
			if got != tt.want {
				t.Errorf("BuildProductURL(%q) = %v, want %v", tt.productID, got, tt.want)
			}
		})
	}
}

func TestEbayScraper_Platform(t *testing.T) {
	s := &EbayScraper{}
	if s.Platform() != PlatformEbay {
		t.Errorf("Platform() = %v, want %v", s.Platform(), PlatformEbay)
	}
}

func TestEbayScraper_Interface(t *testing.T) {
	var s PlatformScraper = &EbayScraper{}
	if s.Platform() != PlatformEbay {
		t.Errorf("Platform() = %v, want %v", s.Platform(), PlatformEbay)
	}
	if s.BuildSearchURL("test") == "" {
		t.Error("BuildSearchURL returned empty string")
	}
	if s.BuildProductURL("123") == "" {
		t.Error("BuildProductURL returned empty string")
	}
}

func TestEbayScraper_MockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>eBay Test</title></head>
<body>
<div class="s-item">
  <a class="s-item__link" href="/itm/123456789">
    <h3 class="s-item__title">Test Product</h3>
  </a>
  <span class="s-item__price">$19.99</span>
  <span class="s-item__subtitle">Used - Good condition</span>
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

func TestEbayScraper_URLBuilding(t *testing.T) {
	s := &EbayScraper{}

	// Test query with special characters
	url := s.BuildSearchURL("laptop")
	expected := "https://www.ebay.com/sch/i.html?_nkw=laptop"
	if url != expected {
		t.Errorf("BuildSearchURL() = %v, want %v", url, expected)
	}

	// Test product URL
	productURL := s.BuildProductURL("384725910")
	expectedProduct := "https://www.ebay.com/itm/384725910"
	if productURL != expectedProduct {
		t.Errorf("BuildProductURL() = %v, want %v", productURL, expectedProduct)
	}
}
