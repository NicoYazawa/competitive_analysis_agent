package api

import (
	"net/http"

	"competitive-analysis-agent/internal/api/handlers"
	apimiddleware "competitive-analysis-agent/internal/api/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	competitorHandler *handlers.CompetitorHandler,
	trendHandler *handlers.TrendHandler,
	pricingHandler *handlers.PricingHandler,
	productHandler *handlers.ProductHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// Setup structured logging and trace ID middleware
	r.Use(apimiddleware.TraceIDMiddleware)
	r.Use(apimiddleware.RequestLogger(apimiddleware.DefaultLogger()))

	// Health check
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Competitor routes
		r.Route("/competitors", func(r chi.Router) {
			r.Get("/", competitorHandler.List)
			r.Get("/search", competitorHandler.Search)
			r.Get("/price-changes", competitorHandler.DetectPriceChanges)
			r.Get("/{id}", competitorHandler.Get)
			r.Get("/{id}/price-history", competitorHandler.GetPriceHistory)
		})

		// Trend routes
		r.Route("/trends", func(r chi.Router) {
			r.Get("/analyze", trendHandler.AnalyzeTrend)
			r.Get("/analysis", trendHandler.GetAggregatedAnalysis)
			r.Get("/categories", trendHandler.ListTrendCategories)
			r.Get("/categories/{category}", trendHandler.GetTrendByCategory)
			r.Get("/demand-signals", trendHandler.GetDemandSignals)
			r.Get("/overview", trendHandler.GetMarketOverview)
			r.Get("/trending", trendHandler.GetTrendingProducts)
		})

		// Strategy/Pricing routes
		r.Route("/strategy", func(r chi.Router) {
			r.Route("/pricing", func(r chi.Router) {
				r.Get("/", pricingHandler.GetPricingRecommendation)
				r.Get("/competitors/{product_id}", pricingHandler.GetCompetitorPricing)
				r.Get("/elasticity/{product_id}", pricingHandler.GetPriceElasticity)
				r.Post("/scenario", pricingHandler.GetPricingScenario)
				r.Get("/history/{product_id}", pricingHandler.GetPricingHistory)
				r.Get("/dynamic/{product_id}", pricingHandler.GetDynamicPricing)
			})
		})

		// Product routes
		r.Route("/products", func(r chi.Router) {
			r.Get("/recommendations", productHandler.GetProductRecommendations)
			r.Get("/compare", productHandler.GetProductComparison)
			r.Get("/trends", productHandler.GetProductTrends)
			r.Get("/{id}", productHandler.GetProductByID)
			r.Get("/{id}/analysis", productHandler.GetProductAnalysis)
			r.Get("/categories/{category}/insights", productHandler.GetCategoryInsights)
		})
	})

	return r
}
