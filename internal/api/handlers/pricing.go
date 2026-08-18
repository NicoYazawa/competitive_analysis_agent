package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"competitive-analysis-agent/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

// PricingHandler handles pricing-related HTTP requests.
type PricingHandler struct{}

// NewPricingHandler creates a new PricingHandler.
func NewPricingHandler() *PricingHandler {
	return &PricingHandler{}
}

// PricingRecommendation represents a pricing recommendation.
type PricingRecommendation struct {
	ProductName       string  `json:"product_name"`
	CurrentPrice      float64 `json:"current_price"`
	RecommendedPrice  float64 `json:"recommended_price"`
	MinPrice          float64 `json:"min_price"`
	MaxPrice          float64 `json:"max_price"`
	Currency          string  `json:"currency"`
	Confidence        float64 `json:"confidence"`
	Positioning       string  `json:"positioning"`
	CompetitorAvg     float64 `json:"competitor_average"`
	MarketDemand      string  `json:"market_demand"`
}

// PricingStrategyResponse represents a complete pricing strategy response.
type PricingStrategyResponse struct {
	ProductID         string                  `json:"product_id"`
	Recommendations   []PricingRecommendation `json:"recommendations"`
	Strategy          string                  `json:"strategy"`
	Summary           string                  `json:"summary"`
}

// GetPricingRecommendation handles GET /api/v1/strategy/pricing?product_id=xxx
func (h *PricingHandler) GetPricingRecommendation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		writeError(w, http.StatusBadRequest, "product_id parameter required", traceID)
		return
	}

	// Generate pricing recommendation based on mock data
	recommendation := PricingRecommendation{
		ProductName:       "Sample Product",
		CurrentPrice:      99.99,
		RecommendedPrice:  94.99,
		MinPrice:         79.99,
		MaxPrice:         129.99,
		Currency:         "USD",
		Confidence:       0.85,
		Positioning:      "Mid-market",
		CompetitorAvg:    102.50,
		MarketDemand:     "High",
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"product_id":      productID,
		"recommendation":   recommendation,
		"trace_id":        traceID,
	}, traceID)
}

// GetCompetitorPricing handles GET /api/v1/strategy/pricing/competitors/{product_id}
func (h *PricingHandler) GetCompetitorPricing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	productID := chi.URLParam(r, "product_id")
	if productID == "" {
		writeError(w, http.StatusBadRequest, "product_id parameter required", traceID)
		return
	}

	type competitorPrice struct {
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Currency    string  `json:"currency"`
		Rating      float64 `json:"rating"`
		ReviewCount int     `json:"review_count"`
		Timestamp  string  `json:"timestamp"`
	}

	competitors := []competitorPrice{
		{Name: "Competitor A", Price: 105.00, Currency: "USD", Rating: 4.5, ReviewCount: 1200, Timestamp: "2026-08-18T10:00:00Z"},
		{Name: "Competitor B", Price: 98.50, Currency: "USD", Rating: 4.2, ReviewCount: 850, Timestamp: "2026-08-18T09:30:00Z"},
		{Name: "Competitor C", Price: 112.00, Currency: "USD", Rating: 4.7, ReviewCount: 2100, Timestamp: "2026-08-18T10:15:00Z"},
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"product_id":  productID,
		"competitors": competitors,
		"avg_price":   105.17,
		"min_price":   98.50,
		"max_price":   112.00,
		"trace_id":    traceID,
	}, traceID)
}

// GetPriceElasticity handles GET /api/v1/strategy/pricing/elasticity/{product_id}
func (h *PricingHandler) GetPriceElasticity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	productID := chi.URLParam(r, "product_id")

	type elasticityResponse struct {
		ProductID    string  `json:"product_id"`
		Elasticity   float64 `json:"elasticity"`
		ElasticityType string `json:"elasticity_type"` // "elastic", "inelastic", "unitary"
		OptimalRange string  `json:"optimal_range"`
		Recommendation string `json:"recommendation"`
		TraceID      string  `json:"trace_id,omitempty"`
	}

	elasticity := elasticityResponse{
		ProductID:     productID,
		Elasticity:    -1.25,
		ElasticityType: "elastic",
		OptimalRange:  "$85 - $105",
		Recommendation: "Price decrease will likely increase revenue",
	}

	writeJSON(w, http.StatusOK, elasticity, traceID)
}

// GetPricingScenario handles POST /api/v1/strategy/pricing/scenario
func (h *PricingHandler) GetPricingScenario(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	if r.Header.Get("Content-Type") != "application/json" {
		writeError(w, http.StatusBadRequest, "Content-Type must be application/json", traceID)
		return
	}

	var req struct {
		ProductID    string  `json:"product_id"`
		BasePrice    float64 `json:"base_price"`
		TargetMargin float64 `json:"target_margin"`
		CompetitorPrices []float64 `json:"competitor_prices"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", traceID)
		return
	}

	type scenarioResult struct {
		Scenario        string  `json:"scenario"`
		Revenue         float64 `json:"projected_revenue"`
		Profit          float64 `json:"projected_profit"`
		Margin          float64 `json:"margin_percent"`
		MarketShareGain float64 `json:"market_share_gain_percent"`
	}

	scenarios := []scenarioResult{
		{Scenario: "Aggressive (-10%)", Revenue: req.BasePrice * 0.9 * 150, Profit: req.BasePrice * 0.9 * 150 * 0.25, Margin: 25, MarketShareGain: 15},
		{Scenario: "Moderate (-5%)", Revenue: req.BasePrice * 0.95 * 130, Profit: req.BasePrice * 0.95 * 130 * 0.30, Margin: 30, MarketShareGain: 8},
		{Scenario: "Maintain", Revenue: req.BasePrice * 110, Profit: req.BasePrice * 110 * 0.35, Margin: 35, MarketShareGain: 0},
		{Scenario: "Premium (+5%)", Revenue: req.BasePrice * 1.05 * 95, Profit: req.BasePrice * 1.05 * 95 * 0.40, Margin: 40, MarketShareGain: -5},
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"product_id": req.ProductID,
		"scenarios":  scenarios,
		"recommended": "Moderate (-5%)",
		"trace_id":  traceID,
	}, traceID)
}

// GetPricingHistory handles GET /api/v1/strategy/pricing/history/{product_id}
func (h *PricingHandler) GetPricingHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	productID := chi.URLParam(r, "product_id")
	daysStr := r.URL.Query().Get("days")
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 365 {
			days = d
		}
	}

	type pricePoint struct {
		Date       string  `json:"date"`
		Price      float64 `json:"price"`
		CompetitorAvg float64 `json:"competitor_avg"`
	}

	points := generateMockPriceHistory(days)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"product_id": productID,
		"days":       days,
		"history":    points,
		"trace_id":  traceID,
	}, traceID)
}

func generateMockPriceHistory(days int) []map[string]interface{} {
	points := make([]map[string]interface{}, days)
	basePrice := 99.99

	for i := 0; i < days; i++ {
		points[i] = map[string]interface{}{
			"date":            "2026-08-" + strconv.Itoa(18-days+i+1),
			"price":           basePrice + float64(i%7)*2 - 6,
			"competitor_avg":  basePrice + 5 + float64(i%5)*1.5,
		}
	}
	return points
}

// GetDynamicPricing handles GET /api/v1/strategy/pricing/dynamic/{product_id}
func (h *PricingHandler) GetDynamicPricing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	productID := chi.URLParam(r, "product_id")

	type dynamicPricing struct {
		ProductID       string  `json:"product_id"`
		CurrentPrice    float64 `json:"current_price"`
		RecommendedPrice float64 `json:"recommended_price"`
		Factors         map[string]string `json:"factors"`
		NextUpdate      string  `json:"next_update"`
		TraceID         string  `json:"trace_id,omitempty"`
	}

	factors := map[string]string{
		"demand":     "High - peak shopping hours",
		"inventory":   "Medium stock levels",
		"competition": "Stable competitor pricing",
		"seasonality": "Back-to-school boost",
	}

	dp := dynamicPricing{
		ProductID:        productID,
		CurrentPrice:     99.99,
		RecommendedPrice: 94.99,
		Factors:         factors,
		NextUpdate:      "2026-08-18T12:00:00Z",
	}

	writeJSON(w, http.StatusOK, dp, traceID)
}
