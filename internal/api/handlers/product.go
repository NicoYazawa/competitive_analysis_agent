package handlers

import (
	"net/http"
	"strconv"

	"competitive-analysis-agent/internal/api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ProductHandler handles product selection and recommendation HTTP requests.
type ProductHandler struct{}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

// ProductRecommendation represents a single product recommendation.
type ProductRecommendation struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Category          string  `json:"category"`
	Score             float64 `json:"score"`
	PriceRange        string  `json:"price_range"`
	DemandLevel       string  `json:"demand_level"`
	CompetitionLevel  string  `json:"competition_level"`
	Recommendation    string  `json:"recommendation"`
	Pros              []string `json:"pros"`
	Cons              []string `json:"cons"`
}

// ProductSelectionResponse represents the response for product selection.
type ProductSelectionResponse struct {
	Products      []ProductRecommendation `json:"products"`
	Total         int                    `json:"total"`
	FilterApplied map[string]string      `json:"filter_applied,omitempty"`
}

// GetProductRecommendations handles GET /api/v1/products/recommendations
func (h *ProductHandler) GetProductRecommendations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	category := r.URL.Query().Get("category")
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = l
		}
	}

	products := []ProductRecommendation{
		{
			ID:               uuid.New().String(),
			Name:             "Wireless Earbuds Pro",
			Category:         "Electronics",
			Score:            9.2,
			PriceRange:       "$50-$150",
			DemandLevel:      "High",
			CompetitionLevel: "Very High",
			Recommendation:   "Strong Buy",
			Pros:             []string{"High demand", "Recurring revenue potential", "Compact"},
			Cons:             []string{"Saturated market", "Technical support required"},
		},
		{
			ID:               uuid.New().String(),
			Name:             "Smart Fitness Tracker",
			Category:         "Electronics",
			Score:            8.8,
			PriceRange:       "$30-$100",
			DemandLevel:      "High",
			CompetitionLevel: "High",
			Recommendation:   "Buy",
			Pros:             []string{"Growing market", "Health trend", "Low return rate"},
			Cons:             []string{"Feature parity competition", "Brand loyalty"},
		},
		{
			ID:               uuid.New().String(),
			Name:             "Ergonomic Office Chair",
			Category:         "Home",
			Score:            8.5,
			PriceRange:       "$150-$500",
			DemandLevel:      "Medium",
			CompetitionLevel: "Medium",
			Recommendation:   "Hold",
			Pros:             []string{"High margin", "Less price sensitive", "Repeat customers"},
			Cons:             []string{"Shipping costs", "Assembly required"},
		},
		{
			ID:               uuid.New().String(),
			Name:             "Portable Power Bank",
			Category:         "Electronics",
			Score:            7.9,
			PriceRange:       "$20-$60",
			DemandLevel:      "High",
			CompetitionLevel: "Very High",
			Recommendation:   "Selective",
			Pros:             []string{"Universal compatibility", "Travel essential"},
			Cons:             []string{"Price wars", "Capacity claims vary"},
		},
		{
			ID:               uuid.New().String(),
			Name:             "Premium Phone Case",
			Category:         "Accessories",
			Score:            7.5,
			PriceRange:       "$15-$40",
			DemandLevel:      "Medium",
			CompetitionLevel: "High",
			Recommendation:   "Selective",
			Pros:             []string{"High margin", "Impulse buy", "Bundle potential"},
			Cons:             []string{"Model-specific", "Trend-dependent"},
		},
	}

	if category != "" {
		filtered := make([]ProductRecommendation, 0)
		for _, p := range products {
			if p.Category == category {
				filtered = append(filtered, p)
			}
		}
		products = filtered
	}

	if limit < len(products) {
		products = products[:limit]
	}

	response := ProductSelectionResponse{
		Products: products,
		Total:    len(products),
		FilterApplied: make(map[string]string),
	}
	if category != "" {
		response.FilterApplied["category"] = category
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"recommendations": response.Products,
		"total":          response.Total,
		"filter_applied": response.FilterApplied,
		"trace_id":       traceID,
	}, traceID)
}

// GetProductByID handles GET /api/v1/products/{id}
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	idStr := chi.URLParam(r, "id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid product ID format", traceID)
		return
	}

	product := map[string]interface{}{
		"id":           idStr,
		"name":         "Sample Product",
		"category":     "Electronics",
		"brand":        "Sample Brand",
		"description":  "This is a sample product description",
		"price":        99.99,
		"currency":     "USD",
		"rating":       4.5,
		"review_count": 1250,
		"demand_score": 8.7,
		"competition":  "Medium",
	}

	writeJSON(w, http.StatusOK, product, traceID)
}

// GetProductAnalysis handles GET /api/v1/products/{id}/analysis
func (h *ProductHandler) GetProductAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	idStr := chi.URLParam(r, "id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid product ID format", traceID)
		return
	}

	type swot struct {
		Strengths     []string `json:"strengths"`
		Weaknesses     []string `json:"weaknesses"`
		Opportunities []string `json:"opportunities"`
		Threats       []string `json:"threats"`
	}

	analysis := map[string]interface{}{
		"product_id":   idStr,
		"swot": swot{
			Strengths:     []string{"High brand recognition", "Quality manufacturing", "Strong reviews"},
			Weaknesses:    []string{"Premium pricing", "Limited distribution"},
			Opportunities: []string{"Emerging markets", "New product variants"},
			Threats:       []string{"New competitors", "Economic downturn"},
		},
		"market_position": "Premium",
		"target_audience": "Tech-savvy professionals",
		"seasonality":     "Q4 peak (Holiday season)",
	}

	writeJSON(w, http.StatusOK, analysis, traceID)
}

// GetProductComparison handles GET /api/v1/products/compare?ids=id1,id2,id3
func (h *ProductHandler) GetProductComparison(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		writeError(w, http.StatusBadRequest, "ids parameter required", traceID)
		return
	}

	type comparisonItem struct {
		ID            string  `json:"id"`
		Name          string  `json:"name"`
		Price         float64 `json:"price"`
		Rating        float64 `json:"rating"`
		DemandScore   float64 `json:"demand_score"`
		Competition   string  `json:"competition_level"`
		Recommendation string `json:"recommendation"`
	}

	items := []comparisonItem{
		{ID: "1", Name: "Product A", Price: 99.99, Rating: 4.5, DemandScore: 9.0, Competition: "High", Recommendation: "Buy"},
		{ID: "2", Name: "Product B", Price: 89.99, Rating: 4.2, DemandScore: 8.5, Competition: "Medium", Recommendation: "Buy"},
		{ID: "3", Name: "Product C", Price: 109.99, Rating: 4.7, DemandScore: 8.8, Competition: "High", Recommendation: "Hold"},
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"products": items,
		"best_value": "Product B",
		"highest_rated": "Product C",
		"trace_id": traceID,
	}, traceID)
}

// GetCategoryInsights handles GET /api/v1/products/categories/{category}/insights
func (h *ProductHandler) GetCategoryInsights(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	category := chi.URLParam(r, "category")

	type categoryInsights struct {
		Category         string   `json:"category"`
		TotalProducts    int      `json:"total_products"`
		AvgPrice         float64  `json:"avg_price"`
		AvgRating        float64  `json:"avg_rating"`
		TopSellers       []string `json:"top_sellers"`
		EmergingProducts []string `json:"emerging_products"`
		MarketTrend      string   `json:"market_trend"`
		Opportunity      string   `json:"opportunity"`
	}

	insights := map[string]categoryInsights{
		"electronics": {
			Category:         "Electronics",
			TotalProducts:   1250,
			AvgPrice:        89.99,
			AvgRating:       4.3,
			TopSellers:     []string{"Wireless Earbuds", "Smart Watches", "Tablets"},
			EmergingProducts: []string{"AR Glasses", "Smart Home Hub", "Drone Cameras"},
			MarketTrend:    "Growing +12% YoY",
			Opportunity:    "Smart home devices remain underserved",
		},
		"home": {
			Category:       "Home",
			TotalProducts: 890,
			AvgPrice:       65.50,
			AvgRating:      4.4,
			TopSellers:    []string{"Ergonomic Chairs", "LED Lighting", "Storage Solutions"},
			EmergingProducts: []string{"Air Purifiers", "Smart Thermostats"},
			MarketTrend:   "Stable +5% YoY",
			Opportunity:  "Eco-friendly home products",
		},
	}

	if insight, ok := insights[category]; ok {
		writeJSON(w, http.StatusOK, insight, traceID)
		return
	}

	// Default insight
	defaultInsight := categoryInsights{
		Category:       category,
		TotalProducts: 500,
		AvgPrice:      75.00,
		AvgRating:     4.2,
		TopSellers:    []string{},
		EmergingProducts: []string{},
		MarketTrend:   "Stable",
		Opportunity:   "Market analysis pending",
	}

	writeJSON(w, http.StatusOK, defaultInsight, traceID)
}

// GetProductTrends handles GET /api/v1/products/trends
func (h *ProductHandler) GetProductTrends(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	traceID := middleware.ExtractTraceID(ctx)

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	type trendPoint struct {
		Date        string  `json:"date"`
		ProductCount int    `json:"product_count"`
		AvgPrice    float64 `json:"avg_price"`
		NewEntrants int     `json:"new_entrants"`
	}

	trends := []trendPoint{
		{Date: "2026-07-19", ProductCount: 1200, AvgPrice: 88.50, NewEntrants: 45},
		{Date: "2026-07-26", ProductCount: 1225, AvgPrice: 87.20, NewEntrants: 52},
		{Date: "2026-08-02", ProductCount: 1250, AvgPrice: 89.10, NewEntrants: 48},
		{Date: "2026-08-09", ProductCount: 1280, AvgPrice: 88.75, NewEntrants: 61},
		{Date: "2026-08-16", ProductCount: 1310, AvgPrice: 89.99, NewEntrants: 58},
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"period":  period,
		"trends": trends,
		"summary": map[string]string{
			"total_growth":    "+9.2%",
			"avg_price_change": "+1.7%",
			"new_products":    "264",
		},
		"trace_id": traceID,
	}, traceID)
}
