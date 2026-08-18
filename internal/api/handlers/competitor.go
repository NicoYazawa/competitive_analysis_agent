package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"competitive-analysis-agent/internal/api/middleware"
	"competitive-analysis-agent/internal/domain/entity"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// CompetitorRepository interface for competitor data access
type CompetitorRepository interface {
	Create(ctx context.Context, c *entity.Competitor) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Competitor, error)
	GetByPlatformAndProductID(ctx context.Context, platform, productID string) (*entity.Competitor, error)
	Update(ctx context.Context, c *entity.Competitor) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByPlatform(ctx context.Context, platform string, limit, offset int) ([]*entity.Competitor, error)
	ListAll(ctx context.Context, limit, offset int) ([]*entity.Competitor, error)
	CountByPlatform(ctx context.Context, platform string) (int, error)
	UpdatePrice(ctx context.Context, id uuid.UUID, price float64, currency string) error
	Upsert(ctx context.Context, c *entity.Competitor) error
}

// PriceHistoryRepository interface for price history data access
type PriceHistoryRepository interface {
	Create(ctx context.Context, p *entity.PriceHistory) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.PriceHistory, error)
	GetByCompetitorID(ctx context.Context, competitorID uuid.UUID, limit int) ([]*entity.PriceHistory, error)
	GetLatest(ctx context.Context, competitorID uuid.UUID) (*entity.PriceHistory, error)
	GetLatestPrices(ctx context.Context, competitorIDs []uuid.UUID) (map[uuid.UUID]*entity.PriceHistory, error)
	GetPriceRange(ctx context.Context, competitorID uuid.UUID, startTime, endTime time.Time) ([]*entity.PriceHistory, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByCompetitorID(ctx context.Context, competitorID uuid.UUID) error
	Count(ctx context.Context, competitorID uuid.UUID) (int, error)
	DetectPriceChange(ctx context.Context, competitorID uuid.UUID, threshold float64) (interface{}, error)
	GetAveragePrice(ctx context.Context, competitorID uuid.UUID, days int) (float64, error)
	GetMinMaxPrice(ctx context.Context, competitorID uuid.UUID, days int) (min, max float64, err error)
}

// CompetitorHandler handles competitor-related HTTP requests.
type CompetitorHandler struct {
	repo      CompetitorRepository
	priceRepo PriceHistoryRepository
}

// NewCompetitorHandler creates a new CompetitorHandler.
func NewCompetitorHandler(repo CompetitorRepository, priceRepo PriceHistoryRepository) *CompetitorHandler {
	return &CompetitorHandler{repo: repo, priceRepo: priceRepo}
}

// ListCompetitorsResponse represents the response for listing competitors.
type ListCompetitorsResponse struct {
	Competitors []*CompetitorResponse `json:"competitors"`
	Total       int                  `json:"total"`
	Limit       int                  `json:"limit"`
	Offset      int                  `json:"offset"`
}

// CompetitorResponse represents a single competitor in API responses.
type CompetitorResponse struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	Platform          string        `json:"platform"`
	PlatformProductID string        `json:"platform_product_id"`
	CurrentPrice      float64       `json:"current_price"`
	Currency          string        `json:"currency"`
	Rating            float64       `json:"rating"`
	ReviewCount       int           `json:"review_count"`
	SellerRating      float64       `json:"seller_rating"`
	SellerReviewCount int           `json:"seller_review_count"`
	SourceURL         string        `json:"source_url"`
	CreatedAt         string        `json:"created_at"`
	UpdatedAt         string        `json:"updated_at"`
	PriceHistory      []*PricePoint `json:"price_history,omitempty"`
}

// PricePoint represents a price history data point.
type PricePoint struct {
	Price      float64 `json:"price"`
	Currency   string  `json:"currency"`
	RecordedAt string  `json:"recorded_at"`
}

// PriceChange represents a detected price change.
type PriceChange struct {
	CompetitorID    uuid.UUID `json:"competitor_id"`
	OldPrice        float64   `json:"old_price"`
	NewPrice        float64   `json:"new_price"`
	Change          float64   `json:"change"`
	ChangePercent   float64   `json:"change_percent"`
	Direction       string    `json:"direction"` // "increased", "decreased", "unchanged", "new"
	OldRecordedAt   time.Time `json:"old_recorded_at,omitempty"`
	NewRecordedAt   time.Time `json:"new_recorded_at,omitempty"`
}

// List handles GET /api/v1/competitors
func (h *CompetitorHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	if h.repo == nil {
		writeError(w, http.StatusServiceUnavailable, "Competitor service not available", traceID)
		return
	}

	platform := r.URL.Query().Get("platform")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	var competitors []*entity.Competitor
	var err error

	if platform != "" {
		competitors, err = h.repo.ListByPlatform(ctx, platform, limit, offset)
	} else {
		competitors, err = h.repo.ListAll(ctx, limit, offset)
	}

	if err != nil {
		middleware.DefaultLogger().Error("Failed to list competitors",
			err,
			middleware.TraceID(traceID),
			middleware.Path(r.URL.Path),
		)
		writeError(w, http.StatusInternalServerError, "Failed to retrieve competitors", traceID)
		return
	}

	var total int
	if platform != "" {
		total, _ = h.repo.CountByPlatform(ctx, platform)
	} else {
		total = len(competitors)
	}

	response := ListCompetitorsResponse{
		Competitors: make([]*CompetitorResponse, len(competitors)),
		Total:       total,
		Limit:       limit,
		Offset:      offset,
	}

	for i, c := range competitors {
		response.Competitors[i] = toCompetitorResponse(c)
	}

	writeJSON(w, http.StatusOK, response, traceID)
}

// Get handles GET /api/v1/competitors/{id}
func (h *CompetitorHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	if h.repo == nil {
		writeError(w, http.StatusServiceUnavailable, "Competitor service not available", traceID)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid competitor ID format", traceID)
		return
	}

	competitor, err := h.repo.GetByID(ctx, id)
	if err != nil {
		middleware.DefaultLogger().Error("Failed to get competitor",
			err,
			middleware.TraceID(traceID),
			middleware.Path(r.URL.Path),
		)
		writeError(w, http.StatusInternalServerError, "Failed to retrieve competitor", traceID)
		return
	}

	if competitor == nil {
		writeError(w, http.StatusNotFound, "Competitor not found", traceID)
		return
	}

	response := toCompetitorResponse(competitor)

	// Get price history if requested
	includeHistory := r.URL.Query().Get("include_history") == "true"
	if includeHistory && h.priceRepo != nil {
		history, err := h.priceRepo.GetByCompetitorID(ctx, id, 30)
		if err == nil {
			response.PriceHistory = make([]*PricePoint, len(history))
			for i, ph := range history {
				response.PriceHistory[i] = &PricePoint{
					Price:      ph.Price,
					Currency:   ph.Currency,
					RecordedAt: ph.RecordedAt.Format("2006-01-02T15:04:05Z07:00"),
				}
			}
		}
	}

	writeJSON(w, http.StatusOK, response, traceID)
}

// Search handles GET /api/v1/competitors/search?q=query
func (h *CompetitorHandler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	if h.repo == nil {
		writeError(w, http.StatusServiceUnavailable, "Competitor service not available", traceID)
		return
	}

	query := r.URL.Query().Get("q")
	platform := r.URL.Query().Get("platform")

	if query == "" && platform == "" {
		writeError(w, http.StatusBadRequest, "Search query or platform filter required", traceID)
		return
	}

	competitors, err := h.repo.ListByPlatform(ctx, platform, 50, 0)
	if err != nil {
		middleware.DefaultLogger().Error("Failed to search competitors",
			err,
			middleware.TraceID(traceID),
			middleware.Path(r.URL.Path),
		)
		writeError(w, http.StatusInternalServerError, "Failed to search competitors", traceID)
		return
	}

	response := ListCompetitorsResponse{
		Competitors: make([]*CompetitorResponse, len(competitors)),
		Total:       len(competitors),
		Limit:       50,
		Offset:      0,
	}

	for i, c := range competitors {
		response.Competitors[i] = toCompetitorResponse(c)
	}

	writeJSON(w, http.StatusOK, response, traceID)
}

// GetPriceHistory handles GET /api/v1/competitors/{id}/price-history
func (h *CompetitorHandler) GetPriceHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	if h.repo == nil || h.priceRepo == nil {
		writeError(w, http.StatusServiceUnavailable, "Competitor service not available", traceID)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid competitor ID format", traceID)
		return
	}

	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	history, err := h.priceRepo.GetByCompetitorID(ctx, id, days*10)
	if err != nil {
		middleware.DefaultLogger().Error("Failed to get price history",
			err,
			middleware.TraceID(traceID),
			middleware.Path(r.URL.Path),
		)
		writeError(w, http.StatusInternalServerError, "Failed to retrieve price history", traceID)
		return
	}

	type historyResponse struct {
		CompetitorID string        `json:"competitor_id"`
		Days         int           `json:"days"`
		DataPoints   []*PricePoint `json:"data_points"`
	}

	response := historyResponse{
		CompetitorID: idStr,
		Days:         days,
		DataPoints:   make([]*PricePoint, len(history)),
	}

	for i, ph := range history {
		response.DataPoints[i] = &PricePoint{
			Price:      ph.Price,
			Currency:   ph.Currency,
			RecordedAt: ph.RecordedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	writeJSON(w, http.StatusOK, response, traceID)
}

// DetectPriceChanges handles GET /api/v1/competitors/price-changes
func (h *CompetitorHandler) DetectPriceChanges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	if h.repo == nil || h.priceRepo == nil {
		writeError(w, http.StatusServiceUnavailable, "Competitor service not available", traceID)
		return
	}

	thresholdStr := r.URL.Query().Get("threshold")
	threshold := 0.01 // Default 1% threshold
	if thresholdStr != "" {
		if t, err := strconv.ParseFloat(thresholdStr, 64); err == nil && t > 0 {
			threshold = t
		}
	}

	// Get all competitors
	competitors, err := h.repo.ListAll(ctx, 100, 0)
	if err != nil {
		middleware.DefaultLogger().Error("Failed to list competitors for price change detection",
			err,
			middleware.TraceID(traceID),
			middleware.Path(r.URL.Path),
		)
		writeError(w, http.StatusInternalServerError, "Failed to detect price changes", traceID)
		return
	}

	type priceChangeResponse struct {
		CompetitorID string  `json:"competitor_id"`
		Name         string  `json:"name"`
		OldPrice     float64 `json:"old_price"`
		NewPrice     float64 `json:"new_price"`
		Change       float64 `json:"change"`
		ChangePercent float64 `json:"change_percent"`
		Direction    string  `json:"direction"`
	}

	var changes []priceChangeResponse
	for _, c := range competitors {
		// Use priceRepo's DetectPriceChange if available
		_ = threshold // Used when calling DetectPriceChange
		latest, err := h.priceRepo.GetLatest(ctx, c.ID)
		if err != nil || latest == nil {
			continue
		}

		// Simple price change detection
		if c.CurrentPrice != latest.Price {
			change := latest.Price - c.CurrentPrice
			changePercent := 0.0
			if c.CurrentPrice > 0 {
				changePercent = (change / c.CurrentPrice) * 100
			}

			direction := "unchanged"
			if change > threshold*c.CurrentPrice {
				direction = "increased"
			} else if change < -threshold*c.CurrentPrice {
				direction = "decreased"
			}

			if direction != "unchanged" {
				changes = append(changes, priceChangeResponse{
					CompetitorID: c.ID.String(),
					Name:         c.Name,
					OldPrice:     c.CurrentPrice,
					NewPrice:     latest.Price,
					Change:       change,
					ChangePercent: changePercent,
					Direction:    direction,
				})
			}
		}
	}

	type detectResponse struct {
		Threshold string                 `json:"threshold"`
		Changes   []priceChangeResponse `json:"changes"`
		Count     int                   `json:"count"`
	}

	writeJSON(w, http.StatusOK, detectResponse{
		Threshold: strconv.FormatFloat(threshold*100, 'f', 2, 64) + "%",
		Changes:   changes,
		Count:     len(changes),
	}, traceID)
}

func toCompetitorResponse(c *entity.Competitor) *CompetitorResponse {
	if c == nil {
		return nil
	}
	return &CompetitorResponse{
		ID:                c.ID.String(),
		Name:              c.Name,
		Platform:          c.Platform,
		PlatformProductID: c.PlatformProductID,
		CurrentPrice:      c.CurrentPrice,
		Currency:          c.Currency,
		Rating:            c.Rating,
		ReviewCount:       c.ReviewCount,
		SellerRating:      c.SellerRating,
		SellerReviewCount: c.SellerReviewCount,
		SourceURL:         c.SourceURL,
		CreatedAt:         c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
