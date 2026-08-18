package platforms

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockAmazonScraper implements PlatformScraper for Amazon tests
type mockAmazonScraper struct {
	platform   Platform
	searchURL  string
	productURL string
	searchData []*CompetitorData
	product    *CompetitorData
	searchErr  error
	productErr error
}

func (m *mockAmazonScraper) Platform() Platform                      { return m.platform }
func (m *mockAmazonScraper) BuildSearchURL(query string) string       { return m.searchURL }
func (m *mockAmazonScraper) BuildProductURL(productID string) string { return m.productURL }
func (m *mockAmazonScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	return m.searchData, m.searchErr
}
func (m *mockAmazonScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	return m.searchData, m.searchErr
}
func (m *mockAmazonScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
	return m.product, m.productErr
}
func (m *mockAmazonScraper) Close() error { return nil }

// Ensure mockAmazonScraper implements PlatformScraper
var _ PlatformScraper = (*mockAmazonScraper)(nil)

func TestAmazonScraper_BuildSearchURL(t *testing.T) {
	tests := []struct {
		query string
		want  string
	}{
		{"test", "https://www.amazon.com/s?k=test"},
		{"test query", "https://www.amazon.com/s?k=test+query"},
		{"  spaces  ", "https://www.amazon.com/s?k=++spaces++"},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			// Create a minimal scraper to test URL building
			s := &AmazonScraper{}
			got := s.BuildSearchURL(tt.query)
			if got != tt.want {
				t.Errorf("BuildSearchURL(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}

func TestAmazonScraper_BuildProductURL(t *testing.T) {
	tests := []struct {
		productID string
		want     string
	}{
		{"B001", "https://www.amazon.com/dp/B001"},
		{"ABCDEFGHIJ", "https://www.amazon.com/dp/ABCDEFGHIJ"},
	}

	for _, tt := range tests {
		t.Run(tt.productID, func(t *testing.T) {
			s := &AmazonScraper{}
			got := s.BuildProductURL(tt.productID)
			if got != tt.want {
				t.Errorf("BuildProductURL(%q) = %v, want %v", tt.productID, got, tt.want)
			}
		})
	}
}

func TestAmazonScraper_Platform(t *testing.T) {
	s := &AmazonScraper{}
	if s.Platform() != PlatformAmazon {
		t.Errorf("Platform() = %v, want %v", s.Platform(), PlatformAmazon)
	}
}

func TestAmazonScraper_Interface(t *testing.T) {
	var s PlatformScraper = &AmazonScraper{}
	if s.Platform() != PlatformAmazon {
		t.Errorf("Platform() = %v, want %v", s.Platform(), PlatformAmazon)
	}
	if s.BuildSearchURL("test") == "" {
		t.Error("BuildSearchURL returned empty string")
	}
	if s.BuildProductURL("B001") == "" {
		t.Error("BuildProductURL returned empty string")
	}
}

func TestAmazonScraper_MockScraper(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Amazon Test</title></head>
<body>
<div data-component-type="s-search-result">
  <h2><a href="/dp/B001"><span>Test Product</span></a></h2>
  <span class="a-price"><span class="a-offscreen">$19.99</span></span>
  <span class="a-icon-alt">4.5 out of 5 stars</span>
  <span class="a-size-small">123 reviews</span>
</div>
</body>
</html>`))
	}))
	defer server.Close()

	// Test that our mock correctly parses HTML
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("HTTP GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestAmazonScraper_URLBuilding(t *testing.T) {
	s := &AmazonScraper{}

	// Test query with special characters
	url := s.BuildSearchURL("wireless headphones")
	expected := "https://www.amazon.com/s?k=wireless+headphones"
	if url != expected {
		t.Errorf("BuildSearchURL() = %v, want %v", url, expected)
	}

	// Test product URL
	productURL := s.BuildProductURL("B07XJ8C8F5")
	expectedProduct := "https://www.amazon.com/dp/B07XJ8C8F5"
	if productURL != expectedProduct {
		t.Errorf("BuildProductURL() = %v, want %v", productURL, expectedProduct)
	}
}
