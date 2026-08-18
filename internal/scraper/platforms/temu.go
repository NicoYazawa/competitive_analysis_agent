package platforms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// TemuScraper scrapes Temu using Rod.
type TemuScraper struct {
	browser *rod.Browser
	proxies []string
	index   int
}

// NewTemuScraper creates a new Temu scraper.
func NewTemuScraper(proxies []string) (*TemuScraper, error) {
	browser := rod.New()
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect browser: %w", err)
	}
	return &TemuScraper{
		browser: browser,
		proxies: proxies,
	}, nil
}

// Platform returns the platform identifier.
func (s *TemuScraper) Platform() Platform {
	return PlatformTemu
}

// BuildSearchURL returns the Temu search URL.
func (s *TemuScraper) BuildSearchURL(query string) string {
	return fmt.Sprintf("https://www.temu.com/search?q=%s",
		strings.ReplaceAll(query, " ", "+"))
}

// BuildProductURL returns the Temu product detail URL.
func (s *TemuScraper) BuildProductURL(productID string) string {
	return fmt.Sprintf("https://www.temu.com/item/%s.html", productID)
}

// ScrapeSearch crawls the Temu search results page.
func (s *TemuScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	return s.ScrapeSearchByURL(ctx, s.BuildSearchURL(query))
}

// ScrapeSearchByURL crawls a specific Temu search URL.
func (s *TemuScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	page, err := s.getPage(ctx)
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if err := page.Timeout(30 * time.Second).Navigate(url); err != nil {
		return nil, fmt.Errorf("navigation failed: %w", err)
	}
	page.WaitLoad()

	// Temu uses a web component based structure
	page.Timeout(15 * time.Second).MustElement(".search-result-list")

	products, err := page.Elements(".search-result-item")
	if err != nil || len(products) == 0 {
		// Try alternative selector
		products, err = page.Elements("[class*='goods-item']")
		if err != nil || len(products) == 0 {
			return nil, fmt.Errorf("no search results: %w", err)
		}
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

// ScrapeProduct crawls a single Temu product detail page.
func (s *TemuScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
	page, err := s.getPage(ctx)
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if err := page.Timeout(30 * time.Second).Navigate(productURL); err != nil {
		return nil, fmt.Errorf("navigation failed: %w", err)
	}
	page.WaitLoad()

	data := &CompetitorData{
		Platform:  string(PlatformTemu),
		SourceURL: productURL,
	}

	if title := page.Timeout(5 * time.Second).MustElement(".goods-title"); title != nil {
		data.Name, _ = title.Text()
		data.Name = strings.TrimSpace(data.Name)
	}

	if price := page.Timeout(5 * time.Second).MustElement(".goods-price"); price != nil {
		priceText, _ := price.Text()
		data.Price = strings.ReplaceAll(priceText, "$", "")
		data.Currency = "USD"
	}

	if seller := page.Timeout(5 * time.Second).MustElement(".seller-name"); seller != nil {
		data.SellerName, _ = seller.Text()
	}

	// Extract product ID from URL
	if idx := strings.Index(productURL, "/item/"); idx != -1 {
		idPart := productURL[idx+6:]
		if len(idPart) > 0 {
			data.PlatformProductID = strings.TrimSuffix(idPart, ".html")
		}
	}

	return data, nil
}

func (s *TemuScraper) parseProductCard(product *rod.Element) *CompetitorData {
	data := &CompetitorData{
		Platform: string(PlatformTemu),
		Currency: "USD",
	}

	if titleEl := product.MustElement(".goods-title"); titleEl != nil {
		data.Name, _ = titleEl.Text()
		data.Name = strings.TrimSpace(data.Name)
	}

	if linkEl := product.MustElement("a"); linkEl != nil {
		href, _ := linkEl.Attribute("href")
		if href != nil {
			data.SourceURL = "https://www.temu.com" + *href
			if idx := strings.Index(*href, "/item/"); idx != -1 {
				idPart := (*href)[idx+6:]
				data.PlatformProductID = strings.TrimSuffix(idPart, ".html")
			}
		}
	}

	if priceEl := product.MustElement(".goods-price"); priceEl != nil {
		priceText, _ := priceEl.Text()
		data.Price = strings.ReplaceAll(priceText, "$", "")
	}

	return data
}

func (s *TemuScraper) getPage(ctx context.Context) (*rod.Page, error) {
	page, err := s.browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, err
	}
	page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	})
	return page, nil
}

// Close closes the browser.
func (s *TemuScraper) Close() {
	if s.browser != nil {
		s.browser.Close()
	}
}
