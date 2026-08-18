package platforms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

// ShopifyStore represents a discovered Shopify store.
type ShopifyStore struct {
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Category string `json:"category"`
	URL      string `json:"url"`
}

// ShopifyScraper scrapes Shopify stores using Colly.
type ShopifyScraper struct {
	stores  []string // domains like "store.myshopify.com"
	c       *colly.Collector
}

// NewShopifyScraper creates a new Shopify scraper with optional store list.
func NewShopifyScraper(stores []string) *ShopifyScraper {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
		colly.MaxDepth(2),
	)
	return &ShopifyScraper{
		stores: stores,
		c:      c,
	}
}

// Platform returns the platform identifier.
func (s *ShopifyScraper) Platform() Platform {
	return PlatformShopify
}

// BuildSearchURL is not applicable for Shopify; use ScrapeAllStores instead.
func (s *ShopifyScraper) BuildSearchURL(query string) string {
	return ""
}

// BuildProductURL is not applicable; products are accessed via store URLs.
func (s *ShopifyScraper) BuildProductURL(productID string) string {
	return ""
}

// ScrapeSearch scrapes all configured Shopify stores for products matching query.
// It uses the /products.json endpoint which returns structured product data.
func (s *ShopifyScraper) ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error) {
	var results []*CompetitorData
	queryLower := strings.ToLower(query)

	for _, domain := range s.stores {
		storeCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		products, err := s.scrapeStore(storeCtx, domain, queryLower)
		cancel()
		if err != nil {
			continue
		}
		results = append(results, products...)
	}
	return results, nil
}

// ScrapeSearchByURL is not applicable for Shopify.
func (s *ShopifyScraper) ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error) {
	return nil, fmt.Errorf("Shopify does not support URL-based search")
}

// ScrapeProduct scrapes a single Shopify product via its JSON API.
func (s *ShopifyScraper) ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error) {
	// productURL can be either full URL or just store domain
	domain := productURL
	if strings.Contains(productURL, "/products/") {
		parts := strings.Split(productURL, "/products/")
		domain = parts[0]
	}

	product, err := s.fetchProductJSON(ctx, domain, productURL)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// ScrapeAllStores scrapes all configured stores and returns all products.
func (s *ShopifyScraper) ScrapeAllStores(ctx context.Context) ([]*CompetitorData, error) {
	var results []*CompetitorData
	for _, domain := range s.stores {
		storeCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		products, err := s.scrapeStore(storeCtx, domain, "")
		cancel()
		if err != nil {
			continue
		}
		results = append(results, products...)
	}
	return results, nil
}

// scrapeStore scrapes a single Shopify store's products.json endpoint.
func (s *ShopifyScraper) scrapeStore(ctx context.Context, domain, filterQuery string) ([]*CompetitorData, error) {
	url := fmt.Sprintf("https://%s/products.json?limit=250", domain)

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", domain, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, domain)
	}

	var result struct {
		Products []struct {
			ID            int64    `json:"id"`
			Title        string   `json:"title"`
			BodyHTML     string   `json:"body_html"`
			Vendor       string   `json:"vendor"`
			ProductType  string   `json:"product_type"`
			CreatedAt    string   `json:"created_at"`
			Handle       string   `json:"handle"`
			UpdatedAt    string   `json:"updated_at"`
			PublishedAt  string   `json:"published_at"`
			Variants     []struct {
				Price      string `json:"price"`
				CompareAt  string `json:"compare_at_price"`
			} `json:"variants"`
			Images []struct {
				Src string `json:"src"`
			} `json:"images"`
		} `json:"products"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON from %s: %w", domain, err)
	}

	var results []*CompetitorData
	for _, p := range result.Products {
		// Filter by query if provided
		if filterQuery != "" && !strings.Contains(strings.ToLower(p.Title), filterQuery) {
			continue
		}

		price := ""
		if len(p.Variants) > 0 {
			price = p.Variants[0].Price
		}

		imageURL := ""
		if len(p.Images) > 0 {
			imageURL = p.Images[0].Src
		}

		desc := p.BodyHTML
		if len(desc) > 500 {
			desc = desc[:500]
		}

		results = append(results, &CompetitorData{
			Name:              p.Title,
			Platform:          string(PlatformShopify),
			PlatformProductID: fmt.Sprintf("%d", p.ID),
			Price:             price,
			Currency:          "USD",
			SourceURL:         fmt.Sprintf("https://%s/products/%s", domain, p.Handle),
			SellerName:        p.Vendor,
			Category:          p.ProductType,
			ImageURL:          imageURL,
			Description:       desc,
		})
	}

	return results, nil
}

// fetchProductJSON fetches a single product from a Shopify store.
func (s *ShopifyScraper) fetchProductJSON(ctx context.Context, domain, productURL string) (*CompetitorData, error) {
	var handle string
	if strings.Contains(productURL, "/products/") {
		parts := strings.Split(productURL, "/products/")
		if len(parts) > 1 {
			handle = strings.TrimSuffix(parts[1], ".html")
		}
	}

	if handle == "" {
		return nil, fmt.Errorf("invalid Shopify product URL: %s", productURL)
	}

	url := fmt.Sprintf("https://%s/products/%s.json", domain, handle)
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Product struct {
			ID           int64  `json:"id"`
			Title       string `json:"title"`
			BodyHTML    string `json:"body_html"`
			Vendor      string `json:"vendor"`
			ProductType string `json:"product_type"`
			Handle      string `json:"handle"`
			Variants    []struct {
				Price string `json:"price"`
			} `json:"variants"`
			Images []struct {
				Src string `json:"src"`
			} `json:"images"`
		} `json:"product"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	p := result.Product
	price := ""
	if len(p.Variants) > 0 {
		price = p.Variants[0].Price
	}
	imageURL := ""
	if len(p.Images) > 0 {
		imageURL = p.Images[0].Src
	}

	return &CompetitorData{
		Name:              p.Title,
		Platform:          string(PlatformShopify),
		PlatformProductID: fmt.Sprintf("%d", p.ID),
		Price:            price,
		Currency:         "USD",
		SourceURL:         productURL,
		SellerName:       p.Vendor,
		Category:         p.ProductType,
		ImageURL:         imageURL,
		Description:      p.BodyHTML,
	}, nil
}

// AddStore adds a Shopify store domain to the scraper.
func (s *ShopifyScraper) AddStore(domain string) {
	s.stores = append(s.stores, domain)
}

// DiscoverStores attempts to discover Shopify stores via Google search.
// This is a placeholder - real implementation would use a Google search API
// or scrape Google's search results for "site:myshopify.com" queries.
func (s *ShopifyScraper) DiscoverStores(ctx context.Context, keywords []string) error {
	// TODO: Integrate with Google Search API or SerpAPI
	// For now, stores must be configured manually or via sitemap
	return nil
}
