package entity

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProduct(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	product := Product{
		ID:          id,
		Name:        "Test Product",
		Category:    "Electronics",
		Brand:       "TestBrand",
		Description: "A test product description",
		ImageURL:    "https://example.com/image.jpg",
		SourceURL:   "https://example.com/product",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, id, product.ID)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, "Electronics", product.Category)
	assert.Equal(t, "TestBrand", product.Brand)
	assert.Equal(t, "A test product description", product.Description)
	assert.Equal(t, "https://example.com/image.jpg", product.ImageURL)
	assert.Equal(t, "https://example.com/product", product.SourceURL)
	assert.Equal(t, now, product.CreatedAt)
	assert.Equal(t, now, product.UpdatedAt)
}

func TestProductValidate(t *testing.T) {
	tests := []struct {
		name    string
		product Product
		wantErr error
	}{
		{
			name: "valid product",
			product: Product{
				Name: "Test Product",
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			product: Product{
				Name: "",
			},
			wantErr: ErrEmptyName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProductSetUpdatedAt(t *testing.T) {
	product := Product{Name: "Test"}
	before := time.Now()
	product.SetUpdatedAt()
	after := time.Now()

	assert.True(t, product.UpdatedAt.After(before) || product.UpdatedAt.Equal(before))
	assert.True(t, product.UpdatedAt.Before(after) || product.UpdatedAt.Equal(after))
}

func TestProductJSONSerialization(t *testing.T) {
	id := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)

	product := Product{
		ID:          id,
		Name:        "Test Product",
		Category:    "Electronics",
		Brand:       "TestBrand",
		Description: "A test product",
		ImageURL:    "https://example.com/image.jpg",
		SourceURL:   "https://example.com/product",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	data, err := json.Marshal(product)
	require.NoError(t, err)

	var decoded Product
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, product.ID, decoded.ID)
	assert.Equal(t, product.Name, decoded.Name)
	assert.Equal(t, product.Category, decoded.Category)
}

func TestCompetitor(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	competitor := Competitor{
		ID:                id,
		Name:              "Test Competitor",
		Platform:          "Amazon",
		PlatformProductID: "B001",
		CurrentPrice:      99.99,
		Currency:          "USD",
		Rating:            4.5,
		ReviewCount:       1000,
		SellerRating:      4.8,
		SellerReviewCount: 500,
		SourceURL:         "https://amazon.com/product",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	assert.Equal(t, id, competitor.ID)
	assert.Equal(t, "Test Competitor", competitor.Name)
	assert.Equal(t, "Amazon", competitor.Platform)
	assert.Equal(t, "B001", competitor.PlatformProductID)
	assert.Equal(t, 99.99, competitor.CurrentPrice)
	assert.Equal(t, "USD", competitor.Currency)
	assert.Equal(t, 4.5, competitor.Rating)
	assert.Equal(t, 1000, competitor.ReviewCount)
	assert.Equal(t, 4.8, competitor.SellerRating)
	assert.Equal(t, 500, competitor.SellerReviewCount)
}

func TestCompetitorValidate(t *testing.T) {
	tests := []struct {
		name       string
		competitor Competitor
		wantErr    error
	}{
		{
			name: "valid competitor",
			competitor: Competitor{
				Name:         "Test Competitor",
				Platform:     "Amazon",
				CurrentPrice: 99.99,
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			competitor: Competitor{
				Name:         "",
				Platform:     "Amazon",
				CurrentPrice: 99.99,
			},
			wantErr: ErrEmptyCompetitorName,
		},
		{
			name: "empty platform",
			competitor: Competitor{
				Name:         "Test",
				Platform:     "",
				CurrentPrice: 99.99,
			},
			wantErr: ErrEmptyPlatform,
		},
		{
			name: "negative price",
			competitor: Competitor{
				Name:         "Test",
				Platform:     "Amazon",
				CurrentPrice: -10,
			},
			wantErr: ErrNegativePrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.competitor.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCompetitorSetUpdatedAt(t *testing.T) {
	competitor := Competitor{Name: "Test", Platform: "Amazon"}
	before := time.Now()
	competitor.SetUpdatedAt()
	after := time.Now()

	assert.True(t, competitor.UpdatedAt.After(before) || competitor.UpdatedAt.Equal(before))
	assert.True(t, competitor.UpdatedAt.Before(after) || competitor.UpdatedAt.Equal(after))
}

func TestCompetitorJSONSerialization(t *testing.T) {
	id := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)

	competitor := Competitor{
		ID:                id,
		Name:              "Test Competitor",
		Platform:          "Amazon",
		PlatformProductID: "B001",
		CurrentPrice:      99.99,
		Currency:          "USD",
		Rating:            4.5,
		ReviewCount:       1000,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	data, err := json.Marshal(competitor)
	require.NoError(t, err)

	var decoded Competitor
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, competitor.ID, decoded.ID)
	assert.Equal(t, competitor.CurrentPrice, decoded.CurrentPrice)
}

func TestPriceHistory(t *testing.T) {
	id := uuid.New()
	competitorID := uuid.New()
	now := time.Now()

	priceHistory := PriceHistory{
		ID:           id,
		CompetitorID: competitorID,
		Price:        89.99,
		Currency:     "USD",
		RecordedAt:   now,
	}

	assert.Equal(t, id, priceHistory.ID)
	assert.Equal(t, competitorID, priceHistory.CompetitorID)
	assert.Equal(t, 89.99, priceHistory.Price)
	assert.Equal(t, "USD", priceHistory.Currency)
	assert.Equal(t, now, priceHistory.RecordedAt)
}

func TestPriceHistoryValidate(t *testing.T) {
	competitorID := uuid.New()

	tests := []struct {
		name        string
		priceHistory PriceHistory
		wantErr     error
	}{
		{
			name: "valid price history",
			priceHistory: PriceHistory{
				CompetitorID: competitorID,
				Price:        89.99,
			},
			wantErr: nil,
		},
		{
			name: "zero price",
			priceHistory: PriceHistory{
				CompetitorID: competitorID,
				Price:        0,
			},
			wantErr: ErrZeroPrice,
		},
		{
			name: "negative price",
			priceHistory: PriceHistory{
				CompetitorID: competitorID,
				Price:        -10,
			},
			wantErr: ErrZeroPrice,
		},
		{
			name: "nil competitor ID",
			priceHistory: PriceHistory{
				CompetitorID: uuid.Nil,
				Price:        89.99,
			},
			wantErr: ErrInvalidCompetitorID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.priceHistory.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPriceHistoryJSONSerialization(t *testing.T) {
	id := uuid.New()
	competitorID := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)

	priceHistory := PriceHistory{
		ID:           id,
		CompetitorID: competitorID,
		Price:        89.99,
		Currency:     "USD",
		RecordedAt:   now,
	}

	data, err := json.Marshal(priceHistory)
	require.NoError(t, err)

	var decoded PriceHistory
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, priceHistory.Price, decoded.Price)
	assert.Equal(t, priceHistory.Currency, decoded.Currency)
}

func TestSupplyChainSignal(t *testing.T) {
	id := uuid.New()
	productID := uuid.New()
	now := time.Now()
	rawData := json.RawMessage(`{"source":"test","severity":"high"}`)

	signal := SupplyChainSignal{
		ID:          id,
		ProductID:   productID,
		SignalType:  "price_spike",
		Source:      "internal",
		Severity:    "high",
		Title:       "Price Increase Alert",
		Description: "Price increased by 20%",
		RawData:     rawData,
		CreatedAt:   now,
	}

	assert.Equal(t, id, signal.ID)
	assert.Equal(t, productID, signal.ProductID)
	assert.Equal(t, "price_spike", signal.SignalType)
	assert.Equal(t, "high", signal.Severity)
	assert.NotNil(t, signal.RawData)
}

func TestSupplyChainSignalValidate(t *testing.T) {
	productID := uuid.New()

	tests := []struct {
		name    string
		signal  SupplyChainSignal
		wantErr error
	}{
		{
			name: "valid signal",
			signal: SupplyChainSignal{
				ProductID:  productID,
				SignalType: "price_spike",
				Severity:   "high",
			},
			wantErr: nil,
		},
		{
			name: "empty signal type",
			signal: SupplyChainSignal{
				ProductID:  productID,
				SignalType: "",
				Severity:   "high",
			},
			wantErr: ErrEmptySignalType,
		},
		{
			name: "empty severity",
			signal: SupplyChainSignal{
				ProductID:  productID,
				SignalType: "price_spike",
				Severity:   "",
			},
			wantErr: ErrEmptySeverity,
		},
		{
			name: "nil product ID",
			signal: SupplyChainSignal{
				ProductID:  uuid.Nil,
				SignalType: "price_spike",
				Severity:   "high",
			},
			wantErr: ErrInvalidProductID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.signal.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSupplyChainSignalJSONSerialization(t *testing.T) {
	id := uuid.New()
	productID := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)
	rawData := json.RawMessage(`{"key":"value"}`)

	signal := SupplyChainSignal{
		ID:          id,
		ProductID:   productID,
		SignalType:  "supply_delay",
		Source:      "supplier_api",
		Severity:    "medium",
		Title:       "Supply Delay",
		Description: "Shipment delayed by 3 days",
		RawData:     rawData,
		CreatedAt:   now,
	}

	data, err := json.Marshal(signal)
	require.NoError(t, err)

	var decoded SupplyChainSignal
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, signal.SignalType, decoded.SignalType)
	assert.Equal(t, signal.Severity, decoded.Severity)
}
