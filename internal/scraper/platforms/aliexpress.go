package platforms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// AliExpressScraper scrapes AliExpress using Rod.
type AliExpressScraper struct {
	browser *rod.Browser
	proxies []string
	index   int
}

// NewAliExpressScraper creates a new AliExpress scraper.
func NewAliExpressScraper(proxies []string) (*AliExpressScraper, error) {
	browser := rod.New()
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect browser: %w", err)
	}
	return &AliExpressScraper{
		browser: browser,
		proxies: proxies,
	}, nil
}

// Platform returns the platform identifier.
func (s *AliExpressScraper) Platform() Platform {
	return PlatformAliExpress
}

// BuildSearchURL returns the AliExpress search URL.
func (s *AliExpressScraper) BuildSearchURL(query string) string {
	return fmt.Sprintf("https://www.aliexpress.com/wholesale?SearchText=%s",
		strings.ReplaceAll(query, " ", "+"))
}

// BuildProductURL returns the AliExpress product detail URL.
func (s *AliExpressScraper) BuildProductURL(productID string) string {
	return fmt.Sprintf("https://www.aliexpress.com/item/%s.html", productID)
}

// ScrapeSearch crawls the AliExpress search results page.
func (s *AliExpressScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	return s.ScrapeSearchByURL(ctx, s.BuildSearchURL(query))
}

// ScrapeSearchByURL crawls a specific AliExpress search URL.
func (s *AliExpressScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	page, err := s.newPage()
	if err != nil {
		return nil, err
	}
	defer page.Close()

	if err := page.Timeout(30*time.Second).Navigate(url); err != nil {
		return nil, fmt.Errorf("navigation failed: %w", err)
	}
	page.WaitLoad()

	page.Timeout(15*time.Second).MustElement(".list-item")

	products, err := page.Elements(".list-item")
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

// ScrapeProduct crawls a single AliExpress product detail page.
func (s *AliExpressScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
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
		Platform:  string(PlatformAliExpress),
		SourceURL: productURL,
	}

	if title := page.Timeout(5*time.Second).MustElement("h1.product-title"); title != nil {
		data.Name, _ = title.Text()
		data.Name = strings.TrimSpace(data.Name)
	}

	if price := page.Timeout(5*time.Second).MustElement(".product-price-value"); price != nil {
		priceText, _ := price.Text()
		data.Price = strings.ReplaceAll(priceText, "$", "")
		data.Currency = "USD"
	}

	if orders := page.Timeout(5*time.Second).MustElement(".product-order-count"); orders != nil {
		data.OrdersCount, _ = orders.Text()
	}

	if store := page.Timeout(5*time.Second).MustElement(".store-name"); store != nil {
		data.SellerName, _ = store.Text()
	}

	// Extract product ID from URL
	if idx := strings.Index(productURL, "/item/"); idx != -1 {
		idPart := productURL[idx+6:]
		if dot := strings.Index(idPart, ".html"); dot != -1 && dot > 0 {
			data.PlatformProductID = idPart[:dot]
		}
	}

	return data, nil
}

func (s *AliExpressScraper) parseProductCard(product *rod.Element) *CompetitorData {
	data := &CompetitorData{
		Platform: string(PlatformAliExpress),
		Currency: "USD",
	}

	if titleEl := product.MustElement(".item-title"); titleEl != nil {
		data.Name, _ = titleEl.Text()
		data.Name = strings.TrimSpace(data.Name)
	}

	if linkEl := product.MustElement("a"); linkEl != nil {
		href, _ := linkEl.Attribute("href")
		if href != nil {
			data.SourceURL = *href
			if idx := strings.Index(*href, "/item/"); idx != -1 {
				idPart := (*href)[idx+6:]
				if dot := strings.Index(idPart, ".html"); dot != -1 && dot > 0 {
					data.PlatformProductID = idPart[:dot]
				}
			}
		}
	}

	if priceEl := product.MustElement(".price"); priceEl != nil {
		priceText, _ := priceEl.Text()
		data.Price = strings.ReplaceAll(priceText, "$", "")
	}

	if ordersEl := product.MustElement(".orders-count"); ordersEl != nil {
		data.OrdersCount, _ = ordersEl.Text()
	}

	return data
}

func (s *AliExpressScraper) newPage() (*rod.Page, error) {
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
func (s *AliExpressScraper) Close() {
	if s.browser != nil {
		s.browser.Close()
	}
}
