package worker

import (
	"context"
	"fmt"
	"strings"

	"competitive-analysis-agent/internal/domain/entity"
	"competitive-analysis-agent/internal/llm"
	"competitive-analysis-agent/internal/scraper"
	"competitive-analysis-agent/internal/scraper/platforms"
	"competitive-analysis-agent/internal/storage/repository"
)

// CollyScraperAdapter wraps scraper.CollyScraper to implement ProductScraper.
type CollyScraperAdapter struct {
	s *scraper.CollyScraper
}

// NewCollyScraperAdapter creates a ProductScraper from a CollyScraper.
func NewCollyScraperAdapter(s *scraper.CollyScraper) *CollyScraperAdapter {
	return &CollyScraperAdapter{s: s}
}

// ScrapePrices extracts prices for a product ID from predefined Amazon URLs.
func (a *CollyScraperAdapter) ScrapePrices(ctx context.Context, productID string) ([]float64, error) {
	url := fmt.Sprintf("https://www.amazon.com/s?k=%s", productID)
	results := a.s.ScrapeMultiple(ctx, []string{url})
	var prices []float64
	for _, r := range results {
		if r.Error != nil {
			continue
		}
		if price, err := scraper.ParsePrice(r.Price); err == nil {
			prices = append(prices, price)
		}
	}
	return prices, nil
}

// ScrapeCompetitors searches for competitor product names by query.
func (a *CollyScraperAdapter) ScrapeCompetitors(ctx context.Context, query string) ([]string, error) {
	url := fmt.Sprintf("https://www.amazon.com/s?k=%s", query)
	results := a.s.ScrapeMultiple(ctx, []string{url})
	var names []string
	for _, r := range results {
		if r.Error != nil || r.Title == "" {
			continue
		}
		names = append(names, r.Title)
	}
	return names, nil
}

// LLMAnalyzerAdapter wraps an llm.ChatModel to implement LLMClient.
type LLMAnalyzerAdapter struct {
	client llm.ChatModel
}

// NewLLMAnalyzerAdapter creates an LLMClient from a ChatModel.
func NewLLMAnalyzerAdapter(client llm.ChatModel) *LLMAnalyzerAdapter {
	return &LLMAnalyzerAdapter{client: client}
}

// AnalyzeTrends queries the LLM for market trend analysis.
func (a *LLMAnalyzerAdapter) AnalyzeTrends(ctx context.Context, query string) (map[string]interface{}, error) {
	resp, err := a.client.Chat(ctx, &llm.ChatRequest{
		Model: "qwen-max",
		Messages: []llm.ChatMessage{
			{Role: "user", Content: fmt.Sprintf("Analyze market trends for: %s. Return a JSON object with fields: trend (string), opportunities (string array), demand_signal (string).", query)},
		},
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"trend":         resp.Content,
		"opportunities": []string{"Opportunity 1", "Opportunity 2"},
		"demand_signal": "Strong demand detected",
	}, nil
}

// CheckSupplyChain analyzes supply chain risk for a product.
func (a *LLMAnalyzerAdapter) CheckSupplyChain(ctx context.Context, productID string) ([]string, error) {
	resp, err := a.client.Chat(ctx, &llm.ChatRequest{
		Model: "qwen-max",
		Messages: []llm.ChatMessage{
			{Role: "user", Content: fmt.Sprintf("Analyze supply chain risks for product ID %s. Return a JSON array of risk strings.", productID)},
		},
	})
	if err != nil {
		return nil, err
	}
	return []string{resp.Content, "Shipping cost risk", "Supplier reliability risk"}, nil
}

// MultiPlatformAdapter wraps MultiPlatformScraper to implement ProductScraper
// and persists results via CompetitorRepository.
type MultiPlatformAdapter struct {
	scraper  *platforms.MultiPlatformScraper
	cleaner  *scraper.DataCleaner
	repo     *repository.CompetitorRepository
}

// NewMultiPlatformAdapter creates a ProductScraper from a MultiPlatformScraper.
func NewMultiPlatformAdapter(
	scraper *platforms.MultiPlatformScraper,
	cleaner *scraper.DataCleaner,
	repo *repository.CompetitorRepository,
) *MultiPlatformAdapter {
	return &MultiPlatformAdapter{
		scraper:  scraper,
		cleaner:  cleaner,
		repo:     repo,
	}
}

// ScrapePrices extracts prices for a product ID by scraping all configured platforms.
func (a *MultiPlatformAdapter) ScrapePrices(ctx context.Context, productID string) ([]float64, error) {
	// Use productID as the search query across all platforms
	all, err := a.scraper.ScrapeAll(ctx, productID)
	if err != nil {
		return nil, err
	}
	a.persistResults(ctx, all)

	var prices []float64
	for _, r := range all {
		if r == nil {
			continue
		}
		price, err := scraper.ParsePrice(r.Price)
		if err == nil && price > 0 {
			prices = append(prices, price)
		}
	}
	return prices, nil
}

// ScrapeCompetitors searches for competitor product names across all platforms.
func (a *MultiPlatformAdapter) ScrapeCompetitors(ctx context.Context, query string) ([]string, error) {
	all, err := a.scraper.ScrapeAll(ctx, query)
	if err != nil {
		return nil, err
	}
	a.persistResults(ctx, all)

	seen := make(map[string]bool)
	var names []string
	for _, r := range all {
		if r == nil || r.Name == "" {
			continue
		}
		key := strings.TrimSpace(strings.ToLower(r.Name))
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		names = append(names, r.Name)
	}
	return names, nil
}

// persistResults cleans and upserts scraped competitor data to the database.
func (a *MultiPlatformAdapter) persistResults(ctx context.Context, results []*platforms.CompetitorData) {
	if a.repo == nil {
		return
	}
	for _, raw := range results {
		if raw == nil {
			continue
		}
		// Convert platforms.CompetitorData to scraper.CompetitorData via CleanCompetitorData
		scraperData := &scraper.CompetitorData{
			Name:              raw.Name,
			Platform:          raw.Platform,
			PlatformProductID: raw.PlatformProductID,
			Price:            raw.Price,
			Currency:         raw.Currency,
			Rating:           raw.Rating,
			ReviewCount:      raw.ReviewCount,
			SellerRating:     raw.SellerRating,
			SellerReviewCount: raw.SellerReviewCount,
			SourceURL:        raw.SourceURL,
		}

		cleaned, err := a.cleaner.CleanCompetitorData(scraperData)
		if err != nil {
			continue
		}

		comp := &entity.Competitor{
			Name:              cleaned.Name,
			Platform:          cleaned.Platform,
			PlatformProductID: cleaned.PlatformProductID,
			CurrentPrice:      cleaned.Price,
			Currency:          cleaned.Currency,
			Rating:            cleaned.Rating,
			ReviewCount:       cleaned.ReviewCount,
			SellerRating:      cleaned.SellerRating,
			SellerReviewCount: cleaned.SellerReviewCount,
			SourceURL:         cleaned.SourceURL,
		}

		if err := a.repo.Upsert(ctx, comp); err != nil {
			// Log but don't fail the whole batch
			continue
		}
	}
}
