package handlers

import (
	"context"
	"net/http"
	"strconv"

	"competitive-analysis-agent/internal/api/middleware"
	"competitive-analysis-agent/internal/supervisor"

	"github.com/go-chi/chi/v5"
)

// TrendHandler handles market trend-related HTTP requests.
type TrendHandler struct {
	marketAgent MarketAgentClient
	scheduler   QueryScheduler
	aggregator  ResultAggregator
}

// MarketAgentClient defines the interface for market analysis.
type MarketAgentClient interface {
	Execute(ctx context.Context, task *supervisor.Task) (any, error)
	Name() string
}

// QueryScheduler defines the interface for query scheduling.
type QueryScheduler interface {
	ScheduleAndExecute(ctx context.Context, query string) ([]*supervisor.Task, error)
}

// ResultAggregator defines the interface for result aggregation.
type ResultAggregator interface {
	Aggregate(tasks []*supervisor.Task) *supervisor.AggregatedResult
}

// NewTrendHandler creates a new TrendHandler.
func NewTrendHandler(marketAgent MarketAgentClient, scheduler QueryScheduler, aggregator ResultAggregator) *TrendHandler {
	return &TrendHandler{
		marketAgent: marketAgent,
		scheduler:   scheduler,
		aggregator:  aggregator,
	}
}

// MarketTrendResponse represents the response for market trend analysis.
type MarketTrendResponse struct {
	Trend        string   `json:"trend"`
	Opportunities []string `json:"opportunities"`
	DemandSignal string   `json:"demand_signal"`
	TraceID      string   `json:"trace_id,omitempty"`
}

// demandSignal represents a market demand signal
type demandSignal struct {
	Keyword      string `json:"keyword"`
	SearchVolume string `json:"search_volume"`
	Trend        string `json:"trend"`
	Competition  string `json:"competition"`
}

// AnalyzeTrend handles GET /api/v1/trends/analyze?query=xxx
func (h *TrendHandler) AnalyzeTrend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	query := r.URL.Query().Get("query")
	if query == "" {
		query = "general market trends"
	}

	if h.marketAgent == nil {
		writeError(w, http.StatusServiceUnavailable, "Market trend analysis not available", traceID)
		return
	}

	task := &supervisor.Task{
		Type:  supervisor.TaskTypeMarketTrend,
		Query: query,
	}

	result, err := h.marketAgent.Execute(ctx, task)
	if err != nil {
		middleware.DefaultLogger().Error("Market trend analysis failed",
			err,
			middleware.TraceID(traceID),
			middleware.Path(r.URL.Path),
		)
		writeError(w, http.StatusInternalServerError, "Market trend analysis failed", traceID)
		return
	}

	if trendResult, ok := result.(*supervisor.MarketTrendResult); ok {
		response := MarketTrendResponse{
			Trend:        trendResult.Trend,
			Opportunities: trendResult.Opportunities,
			DemandSignal: trendResult.DemandSignal,
			TraceID:      traceID,
		}
		writeJSON(w, http.StatusOK, response, traceID)
		return
	}

	writeError(w, http.StatusInternalServerError, "Invalid response from market agent", traceID)
}

// GetAggregatedAnalysis handles GET /api/v1/trends/analysis?query=xxx
// This uses the Supervisor's scheduler to decompose query and aggregate results
func (h *TrendHandler) GetAggregatedAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	query := r.URL.Query().Get("query")
	if query == "" {
		writeError(w, http.StatusBadRequest, "Query parameter required", traceID)
		return
	}

	if h.scheduler == nil || h.aggregator == nil {
		writeError(w, http.StatusServiceUnavailable, "Analysis service not available", traceID)
		return
	}

	tasks, err := h.scheduler.ScheduleAndExecute(ctx, query)
	if err != nil {
		middleware.DefaultLogger().Error("Query analysis failed",
			err,
			middleware.TraceID(traceID),
			middleware.Path(r.URL.Path),
		)
		writeError(w, http.StatusInternalServerError, "Query analysis failed", traceID)
		return
	}

	result := h.aggregator.Aggregate(tasks)
	writeJSON(w, http.StatusOK, result, traceID)
}

// ListTrendCategories handles GET /api/v1/trends/categories
func (h *TrendHandler) ListTrendCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	categories := []map[string]string{
		{"id": "electronics", "name": "Electronics", "icon": "📱"},
		{"id": "fashion", "name": "Fashion", "icon": "👗"},
		{"id": "home", "name": "Home & Garden", "icon": "🏠"},
		{"id": "sports", "name": "Sports & Outdoors", "icon": "⚽"},
		{"id": "beauty", "name": "Beauty & Personal Care", "icon": "💄"},
		{"id": "toys", "name": "Toys & Games", "icon": "🎮"},
		{"id": "automotive", "name": "Automotive", "icon": "🚗"},
		{"id": "food", "name": "Food & Beverages", "icon": "🍔"},
	}

	type categoriesResponse struct {
		Categories []map[string]string `json:"categories"`
		TraceID    string             `json:"trace_id,omitempty"`
	}

	writeJSON(w, http.StatusOK, categoriesResponse{
		Categories: categories,
		TraceID:    traceID,
	}, traceID)
}

// GetTrendByCategory handles GET /api/v1/trends/categories/{category}
func (h *TrendHandler) GetTrendByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	category := chi.URLParam(r, "category")
	if category == "" {
		writeError(w, http.StatusBadRequest, "Category parameter required", traceID)
		return
	}

	// Return mock trend data for the category
	type categoryTrend struct {
		Category      string   `json:"category"`
		Trend         string   `json:"trend"`
		DemandLevel   string   `json:"demand_level"`
		GrowthRate    string   `json:"growth_rate"`
		TopProducts   []string `json:"top_products"`
	}

	trends := map[string]categoryTrend{
		"electronics": {
			Category:    "electronics",
			Trend:       "Smart home devices show strong growth",
			DemandLevel: "High",
			GrowthRate:  "+15%",
			TopProducts: []string{"Smart speakers", "Wireless earbuds", "Tablets"},
		},
		"fashion": {
			Category:    "fashion",
			Trend:       "Sustainable fashion continues to gain traction",
			DemandLevel: "Medium",
			GrowthRate:  "+8%",
			TopProducts: []string{"Eco-friendly materials", "Vintage styles", "Customizable items"},
		},
		"home": {
			Category:    "home",
			Trend:       "Home office upgrades remain popular",
			DemandLevel: "High",
			GrowthRate:  "+12%",
			TopProducts: []string{"Ergonomic chairs", "Standing desks", "LED lighting"},
		},
	}

	if trend, ok := trends[category]; ok {
		writeJSON(w, http.StatusOK, trend, traceID)
		return
	}

	// Default response for unknown categories
	writeJSON(w, http.StatusOK, categoryTrend{
		Category:    category,
		Trend:       "General upward trend",
		DemandLevel: "Medium",
		GrowthRate:  "+5%",
		TopProducts: []string{},
	}, traceID)
}

// GetDemandSignals handles GET /api/v1/trends/demand-signals
func (h *TrendHandler) GetDemandSignals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	platform := r.URL.Query().Get("platform")

	type signalsResponse struct {
		Platform  string          `json:"platform,omitempty"`
		Signals   []demandSignal `json:"signals"`
		UpdatedAt string          `json:"updated_at"`
	}

	signals := []demandSignal{
		{Keyword: "wireless earbuds", SearchVolume: "High", Trend: "Rising", Competition: "High"},
		{Keyword: "smart home hub", SearchVolume: "Medium", Trend: "Stable", Competition: "Medium"},
		{Keyword: "fitness tracker", SearchVolume: "High", Trend: "Rising", Competition: "High"},
		{Keyword: "portable charger", SearchVolume: "Very High", Trend: "Rising", Competition: "Very High"},
		{Keyword: "laptop stand", SearchVolume: "Medium", Trend: "Stable", Competition: "Low"},
	}

	if platform != "" {
		signals = filterSignalsByPlatform(signals, platform)
	}

	writeJSON(w, http.StatusOK, signalsResponse{
		Platform:  platform,
		Signals:   signals,
		UpdatedAt: "2026-08-18T10:00:00Z",
	}, traceID)
}

func filterSignalsByPlatform(signals []demandSignal, platform string) []demandSignal {
	// For demonstration, return all signals. In production, filter based on platform.
	return signals
}

// GetMarketOverview handles GET /api/v1/trends/overview
func (h *TrendHandler) GetMarketOverview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	type overviewResponse struct {
		TotalProducts      int     `json:"total_products"`
		ActiveCompetitors  int     `json:"active_competitors"`
		AvgPriceChange     float64 `json:"avg_price_change_percent"`
		MarketSentiment    string  `json:"market_sentiment"`
		TopCategory        string  `json:"top_category"`
		DataPoints         int     `json:"data_points"`
	}

	overview := overviewResponse{
		TotalProducts:     1523,
		ActiveCompetitors: 87,
		AvgPriceChange:    -2.3,
		MarketSentiment:   "Cautiously Optimistic",
		TopCategory:       "Electronics",
		DataPoints:        45230,
	}

	writeJSON(w, http.StatusOK, overview, traceID)
}

// GetTrendingProducts handles GET /api/v1/trends/trending
func (h *TrendHandler) GetTrendingProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	type trendingProduct struct {
		Rank          int     `json:"rank"`
		Name          string  `json:"name"`
		Category      string  `json:"category"`
		PriceChange   float64 `json:"price_change_percent"`
		DemandScore   float64 `json:"demand_score"`
		Recommendation string  `json:"recommendation"`
	}

	products := []trendingProduct{
		{Rank: 1, Name: "Wireless Earbuds Pro", Category: "Electronics", PriceChange: 5.2, DemandScore: 9.5, Recommendation: "Strong Buy"},
		{Rank: 2, Name: "Smart Fitness Watch", Category: "Electronics", PriceChange: -3.1, DemandScore: 9.2, Recommendation: "Buy"},
		{Rank: 3, Name: "Ergonomic Office Chair", Category: "Home", PriceChange: 2.8, DemandScore: 8.8, Recommendation: "Hold"},
		{Rank: 4, Name: "Portable Power Station", Category: "Electronics", PriceChange: 8.5, DemandScore: 8.7, Recommendation: "Strong Buy"},
		{Rank: 5, Name: "Premium Laptop Stand", Category: "Home", PriceChange: -1.2, DemandScore: 8.5, Recommendation: "Buy"},
	}

	if limit < len(products) {
		products = products[:limit]
	}

	type trendingResponse struct {
		Products []trendingProduct `json:"products"`
		TraceID  string           `json:"trace_id,omitempty"`
	}

	writeJSON(w, http.StatusOK, trendingResponse{Products: products, TraceID: traceID}, traceID)
}
