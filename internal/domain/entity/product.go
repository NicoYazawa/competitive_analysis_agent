package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyName     = errors.New("product name cannot be empty")
	ErrInvalidPrice  = errors.New("price must be positive")
)

type Product struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Category    string     `json:"category"`
	Brand       string     `json:"brand"`
	Description string     `json:"description"`
	ImageURL    string     `json:"image_url"`
	SourceURL   string     `json:"source_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrEmptyName
	}
	return nil
}

func (p *Product) SetUpdatedAt() {
	p.UpdatedAt = time.Now()
}
