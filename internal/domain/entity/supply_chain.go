package entity

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptySignalType  = errors.New("signal type cannot be empty")
	ErrEmptySeverity    = errors.New("severity cannot be empty")
	ErrInvalidProductID = errors.New("product ID cannot be empty")
)

type SupplyChainSignal struct {
	ID          uuid.UUID       `json:"id"`
	ProductID   uuid.UUID       `json:"product_id"`
	SignalType  string          `json:"signal_type"`
	Source      string          `json:"source"`
	Severity    string          `json:"severity"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	RawData     json.RawMessage `json:"raw_data"`
	CreatedAt   time.Time       `json:"created_at"`
}

func (s *SupplyChainSignal) Validate() error {
	if s.SignalType == "" {
		return ErrEmptySignalType
	}
	if s.Severity == "" {
		return ErrEmptySeverity
	}
	if s.ProductID == uuid.Nil {
		return ErrInvalidProductID
	}
	return nil
}
