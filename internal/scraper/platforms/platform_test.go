package platforms

import (
	"context"
	"testing"
)

// mockScraper implements PlatformScraper for testing
type mockScraper struct {
	platform   Platform
	searchURL  string
	productURL string
	searchData []*CompetitorData
	product    *CompetitorData
	searchErr  error
	productErr error
}

func (m *mockScraper) Platform() Platform                      { return m.platform }
func (m *mockScraper) BuildSearchURL(query string) string       { return m.searchURL }
func (m *mockScraper) BuildProductURL(productID string) string { return m.productURL }
func (m *mockScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	return m.searchData, m.searchErr
}
func (m *mockScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	return m.searchData, m.searchErr
}
func (m *mockScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
	return m.product, m.productErr
}
func (m *mockScraper) Close() error { return nil }

// Ensure mockScraper implements PlatformScraper
var _ PlatformScraper = (*mockScraper)(nil)

func TestCompetitorData_Validate(t *testing.T) {
	tests := []struct {
		name    string
		data    *CompetitorData
		wantErr bool
	}{
		{
			name: "valid data",
			data: &CompetitorData{
				Name:     "Test Product",
				Platform: "amazon",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			data: &CompetitorData{
				Platform: "amazon",
			},
			wantErr: true,
		},
		{
			name: "missing platform",
			data: &CompetitorData{
				Name: "Test Product",
			},
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    &CompetitorData{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlatform_Constants(t *testing.T) {
	tests := []struct {
		platform Platform
		want     string
	}{
		{PlatformAmazon, "amazon"},
		{PlatformAliExpress, "aliexpress"},
		{PlatformEbay, "ebay"},
		{PlatformTemu, "temu"},
		{PlatformShopify, "shopify"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if string(tt.platform) != tt.want {
				t.Errorf("Platform = %v, want %v", tt.platform, tt.want)
			}
		})
	}
}

func TestPlatformBaseURLs(t *testing.T) {
	urls := PlatformBaseURLs

	if urls[PlatformAmazon] != "https://www.amazon.com" {
		t.Errorf("Amazon URL = %v, want https://www.amazon.com", urls[PlatformAmazon])
	}
	if urls[PlatformAliExpress] != "https://www.aliexpress.com" {
		t.Errorf("AliExpress URL = %v, want https://www.aliexpress.com", urls[PlatformAliExpress])
	}
	if urls[PlatformEbay] != "https://www.ebay.com" {
		t.Errorf("eBay URL = %v, want https://www.ebay.com", urls[PlatformEbay])
	}
	if urls[PlatformTemu] != "https://www.temu.com" {
		t.Errorf("Temu URL = %v, want https://www.temu.com", urls[PlatformTemu])
	}
}

func TestMockScraper_ImplementsInterface(t *testing.T) {
	scraper := &mockScraper{
		platform:   PlatformAmazon,
		searchURL:  "https://www.amazon.com/s?k=test",
		productURL: "https://www.amazon.com/dp/B001",
		searchData: []*CompetitorData{
			{Name: "Product 1", Platform: "amazon"},
			{Name: "Product 2", Platform: "amazon"},
		},
		product: &CompetitorData{
			Name:              "Single Product",
			Platform:          "amazon",
			PlatformProductID: "B001",
			Price:             "19.99",
			Currency:          "USD",
		},
	}

	// Test Platform()
	if scraper.Platform() != PlatformAmazon {
		t.Errorf("Platform() = %v, want %v", scraper.Platform(), PlatformAmazon)
	}

	// Test BuildSearchURL()
	if scraper.BuildSearchURL("test") != scraper.searchURL {
		t.Errorf("BuildSearchURL() = %v, want %v", scraper.BuildSearchURL("test"), scraper.searchURL)
	}

	// Test BuildProductURL()
	if scraper.BuildProductURL("B001") != scraper.productURL {
		t.Errorf("BuildProductURL() = %v, want %v", scraper.BuildProductURL("B001"), scraper.productURL)
	}

	// Test ScrapeSearch()
	ctx := context.Background()
	data, err := scraper.ScrapeSearch(ctx, "test")
	if err != nil {
		t.Errorf("ScrapeSearch() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("ScrapeSearch() returned %d items, want 2", len(data))
	}

	// Test ScrapeProduct()
	product, err := scraper.ScrapeProduct(ctx, "https://www.amazon.com/dp/B001")
	if err != nil {
		t.Errorf("ScrapeProduct() error = %v", err)
	}
	if product.Name != "Single Product" {
		t.Errorf("ScrapeProduct() name = %v, want Single Product", product.Name)
	}
}

func TestMockScraper_ErrorCases(t *testing.T) {
	scraper := &mockScraper{
		platform:   PlatformEbay,
		searchErr:  assertError{"search error"},
		productErr: assertError{"product error"},
	}

	ctx := context.Background()

	_, err := scraper.ScrapeSearch(ctx, "test")
	if err == nil {
		t.Error("ScrapeSearch() expected error, got nil")
	}

	_, err = scraper.ScrapeProduct(ctx, "https://www.ebay.com/itm/123")
	if err == nil {
		t.Error("ScrapeProduct() expected error, got nil")
	}
}

// assertError implements error interface for testing
type assertError struct {
	msg string
}

func (e assertError) Error() string {
	return e.msg
}
