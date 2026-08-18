package platforms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// EbayScraper scrapes eBay using Rod.
type EbayScraper struct {
	browser *rod.Browser
	proxies []string
	index   int
}

// NewEbayScraper creates a new eBay scraper.
func NewEbayScraper(proxies []string) (*EbayScraper, error) {
	browser := rod.New()
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect browser: %w", err)
	}
	return &EbayScraper{
		browser: browser,
		proxies: proxies,
	}, nil
}

// Platform returns the platform identifier.
func (s *EbayScraper) Platform() Platform {
	return PlatformEbay
}

// BuildSearchURL returns the eBay search URL.
func (s *EbayScraper) BuildSearchURL(query string) string {
	return fmt.Sprintf("https://www.ebay.com/sch/i.html?_nkw=%s",
		strings.ReplaceAll(query, " ", "+"))
}

// BuildProductURL returns the eBay product detail URL.
func (s *EbayScraper) BuildProductURL(productID string) string {
	return fmt.Sprintf("https://www.ebay.com/itm/%s", productID)
}

// ScrapeSearch crawls the eBay search results page.
func (s *EbayScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	return s.ScrapeSearchByURL(ctx, s.BuildSearchURL(query))
}

// ScrapeSearchByURL crawls a specific eBay search URL.
func (s *EbayScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	page, err := s.getPage(ctx)
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if err := page.Timeout(30 * time.Second).Navigate(url); err != nil {
		return nil, fmt.Errorf("navigation failed: %w", err)
	}
	page.WaitLoad()

	page.Timeout(15 * time.Second).MustElement(".s-item")

	products, err := page.Elements(".s-item")
	if err != nil || len(products) == 0 {
		return nil, fmt.Errorf("no search results: %w", err)
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

// ScrapeProduct crawls a single eBay product detail page.
func (s *EbayScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
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
		Platform:  string(PlatformEbay),
		SourceURL: productURL,
	}

	if title := page.Timeout(5 * time.Second).MustElement("h1.x-item-title__mainTitle"); title != nil {
		data.Name, _ = title.Text()
		data.Name = strings.TrimSpace(data.Name)
	}

	if price := page.Timeout(5 * time.Second).MustElement(".x-price-primary span"); price != nil {
		priceText, _ := price.Text()
		data.Price = strings.ReplaceAll(priceText, "$", "")
		data.Currency = "USD"
	}

	if condition := page.Timeout(5 * time.Second).MustElement(".x-item-condition__label"); condition != nil {
		data.Condition, _ = condition.Text()
	}

	if sellerRating := page.Timeout(5 * time.Second).MustElement(".x-sellercard__ratings .x-star-rating"); sellerRating != nil {
		data.SellerRating, _ = sellerRating.Text()
	}

	// Extract item ID from URL
	if idx := strings.Index(productURL, "/itm/"); idx != -1 {
		idPart := productURL[idx+5:]
		if len(idPart) > 0 {
			data.PlatformProductID = idPart
		}
	}

	return data, nil
}

func (s *EbayScraper) parseProductCard(product *rod.Element) *CompetitorData {
	data := &CompetitorData{
		Platform: string(PlatformEbay),
		Currency: "USD",
	}

	if titleEl := product.MustElement(".s-item__title"); titleEl != nil {
		data.Name, _ = titleEl.Text()
		data.Name = strings.TrimSpace(data.Name)
	}

	if linkEl := product.MustElement(".s-item__link"); linkEl != nil {
		href, _ := linkEl.Attribute("href")
		if href != nil {
			data.SourceURL = *href
			if idx := strings.Index(*href, "/itm/"); idx != -1 {
				idPart := (*href)[idx+5:]
				if len(idPart) > 0 {
					data.PlatformProductID = idPart
				}
			}
		}
	}

	if priceEl := product.MustElement(".s-item__price"); priceEl != nil {
		priceText, _ := priceEl.Text()
		data.Price = strings.ReplaceAll(priceText, "$", "")
	}

	if conditionEl := product.MustElement(".s-item__subtitle"); conditionEl != nil {
		data.Condition, _ = conditionEl.Text()
		data.Condition = strings.TrimSpace(data.Condition)
	}

	return data
}

func (s *EbayScraper) getPage(ctx context.Context) (*rod.Page, error) {
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
func (s *EbayScraper) Close() {
	if s.browser != nil {
		s.browser.Close()
	}
}
