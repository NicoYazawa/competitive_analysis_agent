package platforms

import (
	"context"
	"fmt"
)

// Platform represents a supported e-commerce platform.
type Platform string

const (
	PlatformAmazon    Platform = "amazon"
	PlatformAliExpress Platform = "aliexpress"
	PlatformEbay       Platform = "ebay"
	PlatformTemu       Platform = "temu"
	PlatformShopify    Platform = "shopify"
)

// PlatformScraper defines the interface for platform-specific scrapers.
type PlatformScraper interface {
	// Platform returns the platform identifier.
	Platform() Platform

	// BuildSearchURL returns the search page URL for a given query.
	BuildSearchURL(query string) string

	// BuildProductURL returns the product detail page URL for a given product ID.
	BuildProductURL(productID string) string

	// ScrapeSearch crawls the search results page and returns raw competitor data.
	ScrapeSearch(ctx context.Context, query string) ([]*CompetitorData, error)

	// ScrapeProduct crawls a single product detail page.
	ScrapeProduct(ctx context.Context, productURL string) (*CompetitorData, error)

	// ScrapeSearchByURL crawls a specific search URL directly.
	ScrapeSearchByURL(ctx context.Context, url string) ([]*CompetitorData, error)
}

// CompetitorData holds raw scraped competitor data before cleaning.
type CompetitorData struct {
	Name              string `json:"name"`
	Platform          string `json:"platform"`
	PlatformProductID string `json:"platform_product_id"`
	Price             string `json:"price"`
	Currency          string `json:"currency"`
	Rating            string `json:"rating"`
	ReviewCount       string `json:"review_count"`
	SellerRating      string `json:"seller_rating"`
	SellerReviewCount string `json:"seller_review_count"`
	SourceURL         string `json:"source_url"`
	OrdersCount       string `json:"orders_count"`    // AliExpress
	Condition         string `json:"condition"`        // eBay
	SellerName        string `json:"seller_name"`      // eBay/AliExpress
	BSR               string `json:"bsr"`              // Amazon Best Sellers Rank
	Category          string `json:"category"`         // product category
	ImageURL          string `json:"image_url"`        // product image
	Description       string `json:"description"`      // product description
}

// Validate checks if the competitor data has the minimum required fields.
func (c *CompetitorData) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.Platform == "" {
		return fmt.Errorf("platform is required")
	}
	return nil
}

// PlatformBaseURLs returns the base URLs for each platform.
var PlatformBaseURLs = map[Platform]string{
	PlatformAmazon:    "https://www.amazon.com",
	PlatformAliExpress: "https://www.aliexpress.com",
	PlatformEbay:       "https://www.ebay.com",
	PlatformTemu:       "https://www.temu.com",
}
