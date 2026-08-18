package platforms

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"competitive-analysis-agent/internal/scraper"
)

// MultiPlatformScraper aggregates multiple platform scrapers.
type MultiPlatformScraper struct {
	scrapers map[Platform]PlatformScraper
	cleaner *scraper.DataCleaner
}

// NewMultiPlatformScraper creates a new multi-platform scraper with all platforms.
func NewMultiPlatformScraper(
	amazon *AmazonScraper,
	aliexpress *AliExpressScraper,
	ebay *EbayScraper,
	temu *TemuScraper,
	shopify *ShopifyScraper,
	cleaner *scraper.DataCleaner,
) *MultiPlatformScraper {
	m := &MultiPlatformScraper{
		scrapers: make(map[Platform]PlatformScraper),
		cleaner:  cleaner,
	}

	if amazon != nil {
		m.scrapers[PlatformAmazon] = amazon
	}
	if aliexpress != nil {
		m.scrapers[PlatformAliExpress] = aliexpress
	}
	if ebay != nil {
		m.scrapers[PlatformEbay] = ebay
	}
	if temu != nil {
		m.scrapers[PlatformTemu] = temu
	}
	if shopify != nil {
		m.scrapers[PlatformShopify] = shopify
	}

	return m
}

// ScrapeAll scrapes all platforms for a given query concurrently.
func (m *MultiPlatformScraper) ScrapeAll(ctx context.Context, query string) ([]*CompetitorData, error) {
	if len(m.scrapers) == 0 {
		return nil, fmt.Errorf("no scrapers configured")
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []*CompetitorData
		errs    []error
	)

	for platform, s := range m.scrapers {
		wg.Add(1)
		go func(p Platform, scraper PlatformScraper) {
			defer wg.Done()

			data, err := scraper.ScrapeSearch(ctx, query)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("[%s] %w", p, err))
				mu.Unlock()
				return
			}

			mu.Lock()
			results = append(results, data...)
			mu.Unlock()
		}(platform, s)
	}

	wg.Wait()

	if len(results) == 0 && len(errs) > 0 {
		return nil, fmt.Errorf("all platforms failed: %s", strings.Join(errsToStrings(errs), "; "))
	}

	// Deduplicate by platform + product ID
	seen := make(map[string]bool)
	deduped := results[:0]
	for _, r := range results {
		key := fmt.Sprintf("%s:%s", r.Platform, r.PlatformProductID)
		if key != ":" && !seen[key] {
			seen[key] = true
			deduped = append(deduped, r)
		}
	}

	return deduped, nil
}

// ScrapePlatform scrapes a single platform.
func (m *MultiPlatformScraper) ScrapePlatform(ctx context.Context, platform Platform, query string) ([]*CompetitorData, error) {
	scraper, ok := m.scrapers[platform]
	if !ok {
		return nil, fmt.Errorf("no scraper for platform: %s", platform)
	}
	return scraper.ScrapeSearch(ctx, query)
}

// ScrapeProduct scrapes a single product from a specific platform.
func (m *MultiPlatformScraper) ScrapeProduct(ctx context.Context, platform Platform, productURL string) (*CompetitorData, error) {
	scraper, ok := m.scrapers[platform]
	if !ok {
		return nil, fmt.Errorf("no scraper for platform: %s", platform)
	}
	return scraper.ScrapeProduct(ctx, productURL)
}

// Platforms returns the list of configured platforms.
func (m *MultiPlatformScraper) Platforms() []Platform {
	platforms := make([]Platform, 0, len(m.scrapers))
	for p := range m.scrapers {
		platforms = append(platforms, p)
	}
	return platforms
}

// Close closes all underlying scrapers.
func (m *MultiPlatformScraper) Close() {
	for _, s := range m.scrapers {
		if c, ok := s.(interface{ Close() }); ok {
			c.Close()
		}
	}
}

// CleanResults cleans and deduplicates scraped data.
func (m *MultiPlatformScraper) CleanResults(results []*CompetitorData) []*CompetitorData {
	if m.cleaner == nil {
		return results
	}

	var cleaned []*CompetitorData
	seen := make(map[string]bool)

	for _, r := range results {
		if r == nil {
			continue
		}
		if err := r.Validate(); err != nil {
			continue
		}

		// Platform normalization
		cleanedPlatform := m.cleaner.CleanPlatform(r.Platform)
		r.Platform = cleanedPlatform

		// Deduplicate
		key := fmt.Sprintf("%s:%s:%s", r.Platform, r.PlatformProductID, r.Name)
		if !seen[key] {
			seen[key] = true
			cleaned = append(cleaned, r)
		}
	}

	return cleaned
}

func errsToStrings(errs []error) []string {
	strs := make([]string, len(errs))
	for i, e := range errs {
		strs[i] = e.Error()
	}
	return strs
}
