package platforms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// AmazonScraper scrapes Amazon using Rod for anti-bot bypass.
type AmazonScraper struct {
	browser *rod.Browser
	proxies []string
	index   int
}

// NewAmazonScraper creates a new Amazon scraper.
func NewAmazonScraper(proxies []string) (*AmazonScraper, error) {
	browser := rod.New()
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect browser: %w", err)
	}
	return &AmazonScraper{
		browser: browser,
		proxies: proxies,
	}, nil
}

// Platform returns the platform identifier.
func (s *AmazonScraper) Platform() Platform {
	return PlatformAmazon
}

// BuildSearchURL returns the Amazon search page URL.
func (s *AmazonScraper) BuildSearchURL(query string) string {
	return fmt.Sprintf("https://www.amazon.com/s?k=%s", strings.ReplaceAll(query, " ", "+"))
}

// BuildProductURL returns the Amazon product detail page URL.
func (s *AmazonScraper) BuildProductURL(productID string) string {
	return fmt.Sprintf("https://www.amazon.com/dp/%s", productID)
}

// ScrapeSearch crawls the Amazon search results page.
func (s *AmazonScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	return s.ScrapeSearchByURL(ctx, s.BuildSearchURL(query))
}

// ScrapeSearchByURL crawls a specific Amazon search URL.
func (s *AmazonScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	page, err := s.newPage()
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if err := page.Timeout(30*time.Second).Navigate(url); err != nil {
		return nil, fmt.Errorf("navigation failed: %w", err)
	}
	page.WaitLoad()

	// Wait for search results
	page.Timeout(15*time.Second).MustElement("[data-component-type='s-search-result']")

	products, err := page.Elements("[data-component-type='s-search-result']")
	if err != nil || len(products) == 0 {
		return nil, fmt.Errorf("no search results found: %w", err)
	}

	var results []*CompetitorData
	for _, product := range products {
		data := s.parseProductCard(product)
		if data != nil && data.Validate() == nil {
			results = append(results, data)
		}
	}
	return results, nil
}

// ScrapeProduct crawls a single Amazon product detail page.
func (s *AmazonScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
	page, err := s.newPage()
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if err := page.Timeout(30*time.Second).Navigate(productURL); err != nil {
		return nil, fmt.Errorf("navigation failed: %w", err)
	}
	page.WaitLoad()

	data := &CompetitorData{
		Platform:  string(PlatformAmazon),
		SourceURL: productURL,
	}

	// Product title
	if title := page.Timeout(5*time.Second).MustElement("#productTitle"); title != nil {
		data.Name, _ = title.Text()
		data.Name = strings.TrimSpace(data.Name)
	}

	// Price
	if price := page.Timeout(5*time.Second).MustElement(".a-price .a-offscreen"); price != nil {
		data.Price, _ = price.Text()
		data.Price = strings.ReplaceAll(data.Price, "$", "")
		data.Currency = "USD"
	}

	// Rating
	if rating := page.Timeout(5*time.Second).MustElement("#acrPopover .a-icon-alt"); rating != nil {
		data.Rating, _ = rating.Text()
		if idx := strings.Index(data.Rating, " "); idx > 0 {
			data.Rating = strings.Split(data.Rating, " ")[0]
		}
	}

	// Review count
	if reviews := page.Timeout(5*time.Second).MustElement("#acrCustomerReviewText"); reviews != nil {
		data.ReviewCount, _ = reviews.Text()
		data.ReviewCount = strings.ReplaceAll(data.ReviewCount, ",", "")
	}

	// Seller
	if seller := page.Timeout(5*time.Second).MustElement("#sellerName"); seller != nil {
		data.SellerName, _ = seller.Text()
	}

	// BSR
	if bsr := page.Timeout(5*time.Second).MustElement("#SalesRank"); bsr != nil {
		data.BSR, _ = bsr.Text()
	}

	// Extract ASIN
	if idx := strings.Index(productURL, "/dp/"); idx != -1 {
		asin := productURL[idx+4:]
		if len(asin) >= 10 {
			data.PlatformProductID = asin[:10]
		}
	}

	return data, nil
}

// parseProductCard extracts data from an Amazon search result card.
func (s *AmazonScraper) parseProductCard(product *rod.Element) *CompetitorData {
	data := &CompetitorData{Platform: string(PlatformAmazon)}

	// Title
	if titleEl := product.MustElement("h2 a span"); titleEl != nil {
		data.Name, _ = titleEl.Text()
		data.Name = strings.TrimSpace(data.Name)
	}

	// URL and ASIN
	if linkEl := product.MustElement("h2 a"); linkEl != nil {
		href, _ := linkEl.Attribute("href")
		if href != nil {
			data.SourceURL = "https://www.amazon.com" + *href
			if idx := strings.Index(*href, "/dp/"); idx != -1 {
				asin := (*href)[idx+4:]
				if len(asin) >= 10 {
					data.PlatformProductID = asin[:10]
				}
			}
		}
	}

	// Price
	if priceEl := product.MustElement(".a-price .a-offscreen"); priceEl != nil {
		data.Price, _ = priceEl.Text()
		data.Price = strings.ReplaceAll(data.Price, "$", "")
		data.Currency = "USD"
	}

	// Rating
	if ratingEl := product.MustElement(".a-icon-alt"); ratingEl != nil {
		data.Rating, _ = ratingEl.Text()
		if idx := strings.Index(data.Rating, " "); idx > 0 {
			data.Rating = strings.Split(data.Rating, " ")[0]
		}
	}

	// Reviews
	if reviewsEl := product.MustElement("span a-size-small"); reviewsEl != nil {
		data.ReviewCount, _ = reviewsEl.Text()
		data.ReviewCount = strings.ReplaceAll(data.ReviewCount, ",", "")
	}

	return data
}

func (s *AmazonScraper) newPage() (*rod.Page, error) {
	page, err := s.browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, err
	}
	page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})
	page.Evaluate(&rod.EvalOptions{JS: `Object.defineProperty(navigator, 'webdriver', {get: () => undefined})`})
	return page, nil
}

// Close closes the browser.
func (s *AmazonScraper) Close() {
	if s.browser != nil {
		s.browser.Close()
	}
}
