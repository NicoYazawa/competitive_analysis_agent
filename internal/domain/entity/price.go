package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrZeroPrice     = errors.New("price must be greater than zero")
	ErrInvalidCompetitorID = errors.New("competitor ID cannot be empty")
)

type PriceHistory struct {
	ID           uuid.UUID `json:"id"`
	CompetitorID uuid.UUID `json:"competitor_id"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency"`
	RecordedAt   time.Time `json:"recorded_at"`
}

func (p *PriceHistory) Validate() error {
	if p.Price <= 0 {
		return ErrZeroPrice
	}
	if p.CompetitorID == uuid.Nil {
		return ErrInvalidCompetitorID
	}
	return nil
}
