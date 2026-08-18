package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyCompetitorName = errors.New("competitor name cannot be empty")
	ErrEmptyPlatform       = errors.New("platform cannot be empty")
	ErrNegativePrice       = errors.New("price cannot be negative")
)

type Competitor struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Platform          string    `json:"platform"`
	PlatformProductID string    `json:"platform_product_id"`
	CurrentPrice      float64   `json:"current_price"`
	Currency          string    `json:"currency"`
	Rating            float64   `json:"rating"`
	ReviewCount       int       `json:"review_count"`
	SellerRating      float64   `json:"seller_rating"`
	SellerReviewCount int       `json:"seller_review_count"`
	SourceURL         string    `json:"source_url"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (c *Competitor) Validate() error {
	if c.Name == "" {
		return ErrEmptyCompetitorName
	}
	if c.Platform == "" {
		return ErrEmptyPlatform
	}
	if c.CurrentPrice < 0 {
		return ErrNegativePrice
	}
	return nil
}

func (c *Competitor) SetUpdatedAt() {
	c.UpdatedAt = time.Now()
}
