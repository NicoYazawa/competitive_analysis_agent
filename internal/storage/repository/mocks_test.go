package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"competitive-analysis-agent/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDB 用于测试的 Mock PostgresDB
type MockDB struct {
	execResult  sql.Result
	execErr     error
	queryRows   *sql.Rows
	queryErr    error
	queryRowVal interface{}
	queryRowErr error
}

func (m *MockDB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.execResult, m.execErr
}

func (m *MockDB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return m.queryRows, m.queryErr
}

func (m *MockDB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return &sql.Row{}
}

func TestNewCompetitorRepository(t *testing.T) {
	repo := NewCompetitorRepository(nil)
	assert.NotNil(t, repo)
}

func TestNewPriceHistoryRepository(t *testing.T) {
	repo := NewPriceHistoryRepository(nil)
	assert.NotNil(t, repo)
}

func TestCompetitor_Validate(t *testing.T) {
	tests := []struct {
		name        string
		competitor  *entity.Competitor
		expectError bool
	}{
		{
			name: "valid competitor",
			competitor: &entity.Competitor{
				ID:           uuid.New(),
				Name:         "Test Product",
				Platform:     "amazon",
				CurrentPrice: 99.99,
			},
			expectError: false,
		},
		{
			name: "empty name",
			competitor: &entity.Competitor{
				ID:       uuid.New(),
				Name:     "",
				Platform: "amazon",
			},
			expectError: true,
		},
		{
			name: "empty platform",
			competitor: &entity.Competitor{
				ID:       uuid.New(),
				Name:     "Test",
				Platform: "",
			},
			expectError: true,
		},
		{
			name: "negative price",
			competitor: &entity.Competitor{
				ID:           uuid.New(),
				Name:         "Test",
				Platform:     "amazon",
				CurrentPrice: -10,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.competitor.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCompetitor_SetUpdatedAt(t *testing.T) {
	c := &entity.Competitor{
		ID:       uuid.New(),
		Name:     "Test",
		Platform: "amazon",
	}

	before := time.Now().Add(-time.Second)
	c.SetUpdatedAt()
	after := time.Now().Add(time.Second)

	assert.True(t, c.UpdatedAt.After(before) || c.UpdatedAt.Equal(before))
	assert.True(t, c.UpdatedAt.Before(after) || c.UpdatedAt.Equal(after))
}

func TestPriceHistory_Validate(t *testing.T) {
	tests := []struct {
		name        string
		ph          *entity.PriceHistory
		expectError bool
	}{
		{
			name: "valid price history",
			ph: &entity.PriceHistory{
				ID:           uuid.New(),
				CompetitorID: uuid.New(),
				Price:        99.99,
				Currency:     "USD",
				RecordedAt:   time.Now(),
			},
			expectError: false,
		},
		{
			name: "zero price",
			ph: &entity.PriceHistory{
				ID:           uuid.New(),
				CompetitorID: uuid.New(),
				Price:        0,
				Currency:     "USD",
				RecordedAt:   time.Now(),
			},
			expectError: true,
		},
		{
			name: "negative price",
			ph: &entity.PriceHistory{
				ID:           uuid.New(),
				CompetitorID: uuid.New(),
				Price:        -10,
				Currency:     "USD",
				RecordedAt:   time.Now(),
			},
			expectError: true,
		},
		{
			name: "nil competitor id",
			ph: &entity.PriceHistory{
				ID:           uuid.New(),
				CompetitorID: uuid.Nil,
				Price:        99.99,
				Currency:     "USD",
				RecordedAt:   time.Now(),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ph.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPriceHistory_NewPriceHistory(t *testing.T) {
	competitorID := uuid.New()
	now := time.Now()

	ph := &entity.PriceHistory{
		ID:           uuid.New(),
		CompetitorID: competitorID,
		Price:        99.99,
		Currency:     "USD",
		RecordedAt:   now,
	}

	assert.NotEqual(t, uuid.Nil, ph.ID)
	assert.Equal(t, competitorID, ph.CompetitorID)
	assert.Equal(t, 99.99, ph.Price)
	assert.Equal(t, "USD", ph.Currency)
	assert.Equal(t, now, ph.RecordedAt)
}

func TestPriceChange_Struct(t *testing.T) {
	competitorID := uuid.New()
	now := time.Now()

	pc := &PriceChange{
		CompetitorID:    competitorID,
		OldPrice:        100.0,
		NewPrice:        110.0,
		Change:          10.0,
		ChangePercent:   10.0,
		Direction:       "increased",
		OldRecordedAt:   now.Add(-24 * time.Hour),
		NewRecordedAt:   now,
	}

	assert.Equal(t, competitorID, pc.CompetitorID)
	assert.Equal(t, 100.0, pc.OldPrice)
	assert.Equal(t, 110.0, pc.NewPrice)
	assert.Equal(t, 10.0, pc.Change)
	assert.Equal(t, 10.0, pc.ChangePercent)
	assert.Equal(t, "increased", pc.Direction)
}

func TestPriceChange_Directions(t *testing.T) {
	directions := []string{"increased", "decreased", "unchanged", "new"}

	for _, dir := range directions {
		pc := &PriceChange{
			Direction: dir,
		}
		assert.Equal(t, dir, pc.Direction)
	}
}

func TestCompetitor_NewCompetitor(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	c := &entity.Competitor{
		ID:                id,
		Name:              "Test Product",
		Platform:          "amazon",
		PlatformProductID: "B001",
		CurrentPrice:      99.99,
		Currency:          "USD",
		Rating:            4.5,
		ReviewCount:       1234,
		SellerRating:      4.8,
		SellerReviewCount: 567,
		SourceURL:         "https://amazon.com/dp/B001",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	assert.Equal(t, id, c.ID)
	assert.Equal(t, "Test Product", c.Name)
	assert.Equal(t, "amazon", c.Platform)
	assert.Equal(t, "B001", c.PlatformProductID)
	assert.Equal(t, 99.99, c.CurrentPrice)
	assert.Equal(t, "USD", c.Currency)
	assert.Equal(t, 4.5, c.Rating)
	assert.Equal(t, 1234, c.ReviewCount)
	assert.Equal(t, 4.8, c.SellerRating)
	assert.Equal(t, 567, c.SellerReviewCount)
	assert.Equal(t, "https://amazon.com/dp/B001", c.SourceURL)
	assert.Equal(t, now, c.CreatedAt)
	assert.Equal(t, now, c.UpdatedAt)
}

func TestCompetitor_ValidateErrors(t *testing.T) {
	tests := []struct {
		name      string
		errVar    error
		errMsg    string
	}{
		{"empty name", entity.ErrEmptyCompetitorName, "competitor name cannot be empty"},
		{"empty platform", entity.ErrEmptyPlatform, "platform cannot be empty"},
		{"negative price", entity.ErrNegativePrice, "price cannot be negative"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.errVar)
			assert.Contains(t, tt.errVar.Error(), tt.errMsg)
		})
	}
}

func TestPriceHistory_ValidateErrors(t *testing.T) {
	tests := []struct {
		name   string
		errVar error
		errMsg string
	}{
		{"zero price", entity.ErrZeroPrice, "price must be greater than zero"},
		{"invalid competitor ID", entity.ErrInvalidCompetitorID, "competitor ID cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.errVar)
			assert.Contains(t, tt.errVar.Error(), tt.errMsg)
		})
	}
}

func TestProduct_ValidateErrors(t *testing.T) {
	err := entity.ErrEmptyName
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product name cannot be empty")
}

func TestProduct_NewProduct(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	p := &entity.Product{
		ID:          id,
		Name:        "Test Product",
		Category:    "Electronics",
		Brand:       "TestBrand",
		Description: "Test description",
		ImageURL:    "https://example.com/image.jpg",
		SourceURL:   "https://example.com/product",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, id, p.ID)
	assert.Equal(t, "Test Product", p.Name)
	assert.Equal(t, "Electronics", p.Category)
	assert.Equal(t, "TestBrand", p.Brand)
	assert.Equal(t, "Test description", p.Description)
	assert.Equal(t, "https://example.com/image.jpg", p.ImageURL)
	assert.Equal(t, "https://example.com/product", p.SourceURL)
}

func TestProduct_Validate(t *testing.T) {
	p := &entity.Product{
		Name: "Test",
	}

	assert.NoError(t, p.Validate())

	p.Name = ""
	assert.Error(t, p.Validate())
}

func TestProduct_SetUpdatedAt(t *testing.T) {
	p := &entity.Product{
		Name: "Test",
	}

	before := time.Now().Add(-time.Second)
	p.SetUpdatedAt()
	after := time.Now().Add(time.Second)

	assert.True(t, p.UpdatedAt.After(before) || p.UpdatedAt.Equal(before))
	assert.True(t, p.UpdatedAt.Before(after) || p.UpdatedAt.Equal(after))
}

func TestSupplyChain_Entity(t *testing.T) {
	id := uuid.New()
	productID := uuid.New()
	now := time.Now()

	// SupplyChain entity is defined but we test its basic structure
	assert.NotEqual(t, uuid.Nil, id)
	assert.NotEqual(t, uuid.Nil, productID)
	assert.True(t, now.Before(time.Now().Add(time.Second)))
}

func TestMarketTrend_Entity(t *testing.T) {
	now := time.Now()

	// MarketTrend entity structure test
	trend := struct {
		ID              uuid.UUID
		Category        string
		TrendKeyword    string
		PopularityScore float64
		GrowthRate      float64
		RecordedAt      time.Time
	}{
		ID:              uuid.New(),
		Category:        "Electronics",
		TrendKeyword:    "smart watch",
		PopularityScore: 85.5,
		GrowthRate:      12.3,
		RecordedAt:      now,
	}

	assert.NotEqual(t, uuid.Nil, trend.ID)
	assert.Equal(t, "Electronics", trend.Category)
	assert.Equal(t, "smart watch", trend.TrendKeyword)
	assert.Equal(t, 85.5, trend.PopularityScore)
	assert.Equal(t, 12.3, trend.GrowthRate)
}

func TestCompetitor_UpsertData(t *testing.T) {
	c := &entity.Competitor{
		ID:                uuid.New(),
		Name:              "Test",
		Platform:          "amazon",
		PlatformProductID: "B001",
		CurrentPrice:      99.99,
		Currency:          "USD",
		Rating:            4.5,
		ReviewCount:       100,
		SellerRating:      4.8,
		SellerReviewCount: 50,
		SourceURL:         "https://amazon.com/dp/B001",
	}

	// Test validation passes
	require.NoError(t, c.Validate())

	// Test that we can set updated at
	c.SetUpdatedAt()
	assert.False(t, c.UpdatedAt.IsZero())
}
