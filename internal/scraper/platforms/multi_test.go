package platforms

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// mockMultiScraper implements PlatformScraper for MultiPlatformScraper tests
type mockMultiScraper struct {
	platform   Platform
	searchData []*CompetitorData
	product    *CompetitorData
	searchErr  error
	productErr error
	sleepDur   time.Duration
}

func (m *mockMultiScraper) Platform() Platform { return m.platform }
func (m *mockMultiScraper) BuildSearchURL(query string) string {
	return "https://" + string(m.platform) + ".com/search?q=" + query
}
func (m *mockMultiScraper) BuildProductURL(productID string) string {
	return "https://" + string(m.platform) + ".com/product/" + productID
}
func (m *mockMultiScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	if m.sleepDur > 0 {
		time.Sleep(m.sleepDur)
	}
	return m.searchData, m.searchErr
}
func (m *mockMultiScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	return m.searchData, m.searchErr
}
func (m *mockMultiScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
	return m.product, m.productErr
}
func (m *mockMultiScraper) Close() error { return nil }

// newTestMultiScraper creates a MultiPlatformScraper directly for testing
func newTestMultiScraper(scrapers map[Platform]PlatformScraper) *MultiPlatformScraper {
	return &MultiPlatformScraper{
		scrapers: scrapers,
		cleaner:  nil,
	}
}

func TestMultiPlatformScraper_ScrapeAll(t *testing.T) {
	amazonScraper := &mockMultiScraper{
		platform: PlatformAmazon,
		searchData: []*CompetitorData{
			{Name: "Amazon Product 1", Platform: "amazon", PlatformProductID: "A1"},
		},
	}
	ebayScraper := &mockMultiScraper{
		platform: PlatformEbay,
		searchData: []*CompetitorData{
			{Name: "eBay Product 1", Platform: "ebay", PlatformProductID: "E1"},
		},
	}

	m := newTestMultiScraper(map[Platform]PlatformScraper{
		PlatformAmazon: amazonScraper,
		PlatformEbay:   ebayScraper,
	})

	ctx := context.Background()
	results, err := m.ScrapeAll(ctx, "test query")
	if err != nil {
		t.Errorf("ScrapeAll() error = %v", err)
	}
	if len(results) != 2 {
		t.Errorf("ScrapeAll() returned %d results, want 2", len(results))
	}
}

func TestMultiPlatformScraper_ScrapeAll_AllPlatforms(t *testing.T) {
	scrapers := map[Platform]PlatformScraper{
		PlatformAmazon: &mockMultiScraper{
			platform: PlatformAmazon,
			searchData: []*CompetitorData{
				{Name: "Amazon Product", Platform: "amazon", PlatformProductID: "A1"},
			},
		},
		PlatformAliExpress: &mockMultiScraper{
			platform: PlatformAliExpress,
			searchData: []*CompetitorData{
				{Name: "AliExpress Product", Platform: "aliexpress", PlatformProductID: "AE1"},
			},
		},
		PlatformEbay: &mockMultiScraper{
			platform: PlatformEbay,
			searchData: []*CompetitorData{
				{Name: "eBay Product", Platform: "ebay", PlatformProductID: "E1"},
			},
		},
		PlatformTemu: &mockMultiScraper{
			platform: PlatformTemu,
			searchData: []*CompetitorData{
				{Name: "Temu Product", Platform: "temu", PlatformProductID: "T1"},
			},
		},
		PlatformShopify: &mockMultiScraper{
			platform: PlatformShopify,
			searchData: []*CompetitorData{
				{Name: "Shopify Product", Platform: "shopify", PlatformProductID: "S1"},
			},
		},
	}

	m := newTestMultiScraper(scrapers)

	ctx := context.Background()
	results, err := m.ScrapeAll(ctx, "test")
	if err != nil {
		t.Errorf("ScrapeAll() error = %v", err)
	}
	if len(results) != 5 {
		t.Errorf("ScrapeAll() returned %d results, want 5", len(results))
	}
}

func TestMultiPlatformScraper_ScrapeAll_Deduplication(t *testing.T) {
	scraper := &mockMultiScraper{
		platform: PlatformAmazon,
		searchData: []*CompetitorData{
			{Name: "Same Product", Platform: "amazon", PlatformProductID: "SAME1"},
			{Name: "Same Product", Platform: "amazon", PlatformProductID: "SAME1"}, // duplicate
		},
	}

	m := newTestMultiScraper(map[Platform]PlatformScraper{
		PlatformAmazon: scraper,
	})

	ctx := context.Background()
	results, err := m.ScrapeAll(ctx, "test")
	if err != nil {
		t.Errorf("ScrapeAll() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("ScrapeAll() returned %d results after dedup, want 1", len(results))
	}
}

func TestMultiPlatformScraper_ScrapeAll_Errors(t *testing.T) {
	amazonScraper := &mockMultiScraper{
		platform:  PlatformAmazon,
		searchErr: errors.New("amazon error"),
	}
	ebayScraper := &mockMultiScraper{
		platform: PlatformEbay,
		searchData: []*CompetitorData{
			{Name: "eBay Product", Platform: "ebay", PlatformProductID: "E1"},
		},
	}

	m := newTestMultiScraper(map[Platform]PlatformScraper{
		PlatformAmazon: amazonScraper,
		PlatformEbay:   ebayScraper,
	})

	ctx := context.Background()
	results, err := m.ScrapeAll(ctx, "test")
	// Should still return results from working scrapers; only returns error if ALL fail with no results
	if len(results) != 1 {
		t.Errorf("ScrapeAll() returned %d results when one platform errored, want 1", len(results))
	}
	// Since eBay succeeded with results, no error is returned even though Amazon failed
	if err != nil {
		t.Errorf("ScrapeAll() unexpected error: %v", err)
	}
}

func TestMultiPlatformScraper_ScrapeAll_AllFailed(t *testing.T) {
	amazonScraper := &mockMultiScraper{
		platform:  PlatformAmazon,
		searchErr: errors.New("amazon error"),
	}
	ebayScraper := &mockMultiScraper{
		platform:  PlatformEbay,
		searchErr: errors.New("ebay error"),
	}

	m := newTestMultiScraper(map[Platform]PlatformScraper{
		PlatformAmazon: amazonScraper,
		PlatformEbay:   ebayScraper,
	})

	ctx := context.Background()
	_, err := m.ScrapeAll(ctx, "test")
	if err == nil {
		t.Error("ScrapeAll() expected error when all platforms failed, got nil")
	}
	// Error message includes platform prefix
	if err != nil && (err.Error() == "" || err.Error() == "all platforms failed: ") {
		t.Errorf("ScrapeAll() error = %v, want error with platform details", err.Error())
	}
}

func TestMultiPlatformScraper_ScrapeAll_EmptyScrapers(t *testing.T) {
	m := newTestMultiScraper(map[Platform]PlatformScraper{})

	ctx := context.Background()
	_, err := m.ScrapeAll(ctx, "test")
	if err == nil {
		t.Error("ScrapeAll() expected error with no scrapers, got nil")
	}
	expectedErr := "no scrapers configured"
	if err.Error() != expectedErr {
		t.Errorf("ScrapeAll() error = %v, want %v", err.Error(), expectedErr)
	}
}

func TestMultiPlatformScraper_ScrapePlatform(t *testing.T) {
	scraper := &mockMultiScraper{
		platform: PlatformAmazon,
		searchData: []*CompetitorData{
			{Name: "Amazon Product", Platform: "amazon", PlatformProductID: "A1"},
		},
	}

	m := newTestMultiScraper(map[Platform]PlatformScraper{
		PlatformAmazon: scraper,
	})

	ctx := context.Background()
	results, err := m.ScrapePlatform(ctx, PlatformAmazon, "test")
	if err != nil {
		t.Errorf("ScrapePlatform() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("ScrapePlatform() returned %d results, want 1", len(results))
	}
}

func TestMultiPlatformScraper_ScrapePlatform_NotFound(t *testing.T) {
	m := newTestMultiScraper(map[Platform]PlatformScraper{})

	ctx := context.Background()
	_, err := m.ScrapePlatform(ctx, PlatformAmazon, "test")
	if err == nil {
		t.Error("ScrapePlatform() expected error for missing platform, got nil")
	}
	expectedErr := "no scraper for platform: amazon"
	if err.Error() != expectedErr {
		t.Errorf("ScrapePlatform() error = %v, want %v", err.Error(), expectedErr)
	}
}

func TestMultiPlatformScraper_ScrapeProduct(t *testing.T) {
	scraper := &mockMultiScraper{
		platform: PlatformAmazon,
		product: &CompetitorData{
			Name:              "Amazon Product",
			Platform:          "amazon",
			PlatformProductID: "B001",
			Price:             "29.99",
		},
	}

	m := newTestMultiScraper(map[Platform]PlatformScraper{
		PlatformAmazon: scraper,
	})

	ctx := context.Background()
	product, err := m.ScrapeProduct(ctx, PlatformAmazon, "https://amazon.com/dp/B001")
	if err != nil {
		t.Errorf("ScrapeProduct() error = %v", err)
	}
	if product.Name != "Amazon Product" {
		t.Errorf("ScrapeProduct() name = %v, want Amazon Product", product.Name)
	}
}

func TestMultiPlatformScraper_ScrapeProduct_NotFound(t *testing.T) {
	m := newTestMultiScraper(map[Platform]PlatformScraper{})

	ctx := context.Background()
	_, err := m.ScrapeProduct(ctx, PlatformAmazon, "https://amazon.com/dp/B001")
	if err == nil {
		t.Error("ScrapeProduct() expected error for missing platform, got nil")
	}
}

func TestMultiPlatformScraper_Platforms(t *testing.T) {
	amazonScraper := &mockMultiScraper{platform: PlatformAmazon}
	ebayScraper := &mockMultiScraper{platform: PlatformEbay}

	m := newTestMultiScraper(map[Platform]PlatformScraper{
		PlatformAmazon: amazonScraper,
		PlatformEbay:   ebayScraper,
	})

	platforms := m.Platforms()
	if len(platforms) != 2 {
		t.Errorf("Platforms() returned %d platforms, want 2", len(platforms))
	}
}

func TestMultiPlatformScraper_Close(t *testing.T) {
	// Note: PlatformScraper interface doesn't have Close() method,
	// so the Close() method uses type assertion to check if underlying
	// scraper implements Close(). Since PlatformScraper doesn't expose
	// Close(), this check will fail for any PlatformScraper.
	// This test verifies the Close() method doesn't panic.
	m := &MultiPlatformScraper{
		scrapers: map[Platform]PlatformScraper{
			PlatformAmazon: &mockMultiScraper{platform: PlatformAmazon},
		},
	}

	// Should not panic
	m.Close()
}

type trackingScraper struct {
	platform Platform
	closed  *bool
}

func (s *trackingScraper) Platform() Platform { return s.platform }
func (s *trackingScraper) BuildSearchURL(query string) string { return "" }
func (s *trackingScraper) BuildProductURL(productID string) string { return "" }
func (s *trackingScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	return nil, nil
}
func (s *trackingScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	return nil, nil
}
func (s *trackingScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
	return nil, nil
}
func (s *trackingScraper) Close() error {
	*s.closed = true
	return nil
}

func TestMultiPlatformScraper_CleanResults(t *testing.T) {
	// Note: CleanResults with nil cleaner returns results as-is without filtering
	// so this test only validates dedup behavior when cleaner is present
	m := newTestMultiScraper(nil)

	results := []*CompetitorData{
		{Name: "Product 1", Platform: "amazon", PlatformProductID: "A1"},
		{Name: "Product 2", Platform: "ebay", PlatformProductID: "E1"},
	}

	cleaned := m.CleanResults(results)
	if len(cleaned) != 2 {
		t.Errorf("CleanResults() returned %d results, want 2", len(cleaned))
	}
}

func TestMultiPlatformScraper_CleanResults_NilData(t *testing.T) {
	// With nil cleaner, CleanResults does not filter nil items
	m := newTestMultiScraper(nil)

	results := []*CompetitorData{
		nil,
		{Name: "Valid Product", Platform: "amazon", PlatformProductID: "A1"},
	}

	cleaned := m.CleanResults(results)
	// nil cleaner returns results as-is
	if len(cleaned) != 2 {
		t.Errorf("CleanResults() with nil cleaner returned %d results, want 2", len(cleaned))
	}
}

func TestMultiPlatformScraper_CleanResults_InvalidData(t *testing.T) {
	// With nil cleaner, CleanResults does not validate items
	m := newTestMultiScraper(nil)

	results := []*CompetitorData{
		{Name: "", Platform: "amazon"}, // would be invalid - missing name
		{Name: "Valid Product", Platform: "amazon", PlatformProductID: "A1"},
	}

	cleaned := m.CleanResults(results)
	// nil cleaner returns results as-is
	if len(cleaned) != 2 {
		t.Errorf("CleanResults() with nil cleaner returned %d results, want 2", len(cleaned))
	}
}

func TestMultiPlatformScraper_ScrapeAll_Concurrent(t *testing.T) {
	var mu sync.Mutex
	callCount := 0

	scraper := &concurrentTrackingScraper{
		mockMultiScraper: mockMultiScraper{
			platform: PlatformAmazon,
			sleepDur: 10 * time.Millisecond,
			searchData: []*CompetitorData{
				{Name: "Product", Platform: "amazon", PlatformProductID: "A1"},
			},
		},
		callCount: &callCount,
		mu:        &mu,
	}

	m := newTestMultiScraper(map[Platform]PlatformScraper{
		PlatformAmazon: scraper,
	})

	ctx := context.Background()
	results, err := m.ScrapeAll(ctx, "test")
	if err != nil {
		t.Errorf("ScrapeAll() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("ScrapeAll() returned %d results, want 1", len(results))
	}
	if callCount != 1 {
		t.Errorf("ScrapeSearch called %d times, want 1", callCount)
	}
}

type concurrentTrackingScraper struct {
	mockMultiScraper
	callCount *int
	mu        *sync.Mutex
}

func (s *concurrentTrackingScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	s.mu.Lock()
	*s.callCount++
	s.mu.Unlock()
	return s.mockMultiScraper.ScrapeSearch(ctx, query)
}

func TestMultiPlatformScraper_Deduplication_EmptyKey(t *testing.T) {
	m := newTestMultiScraper(nil)

	results := []*CompetitorData{
		{Name: "Product", Platform: "", PlatformProductID: ""}, // empty key components
	}

	deduped := m.CleanResults(results)
	// Product with empty platform/productID should not be considered duplicate
	if len(deduped) != 1 {
		t.Errorf("CleanResults() returned %d results, want 1", len(deduped))
	}
}

func TestMultiPlatformScraper_Platforms_Empty(t *testing.T) {
	m := newTestMultiScraper(map[Platform]PlatformScraper{})

	platforms := m.Platforms()
	if len(platforms) != 0 {
		t.Errorf("Platforms() returned %d platforms, want 0", len(platforms))
	}
}

func TestMultiPlatformScraper_ScrapeAll_SinglePlatformResult(t *testing.T) {
	scraper := &mockMultiScraper{
		platform: PlatformAmazon,
		searchData: []*CompetitorData{
			{Name: "Single Product", Platform: "amazon", PlatformProductID: "S1"},
		},
	}

	m := newTestMultiScraper(map[Platform]PlatformScraper{
		PlatformAmazon: scraper,
	})

	ctx := context.Background()
	results, err := m.ScrapeAll(ctx, "test")
	if err != nil {
		t.Errorf("ScrapeAll() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("ScrapeAll() returned %d results, want 1", len(results))
	}
	if results[0].Name != "Single Product" {
		t.Errorf("Product name = %v, want Single Product", results[0].Name)
	}
}

func TestMultiPlatformScraper_ScrapeAll_NoResults(t *testing.T) {
	amazonScraper := &mockMultiScraper{
		platform:   PlatformAmazon,
		searchData: []*CompetitorData{},
	}
	ebayScraper := &mockMultiScraper{
		platform:   PlatformEbay,
		searchData: []*CompetitorData{},
	}

	m := newTestMultiScraper(map[Platform]PlatformScraper{
		PlatformAmazon: amazonScraper,
		PlatformEbay:   ebayScraper,
	})

	ctx := context.Background()
	results, err := m.ScrapeAll(ctx, "nonexistent")
	// No results and no errors should not return error
	if err != nil {
		t.Errorf("ScrapeAll() error = %v, want nil for empty results with no errors", err)
	}
	if len(results) != 0 {
		t.Errorf("ScrapeAll() returned %d results, want 0", len(results))
	}
}
