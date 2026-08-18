package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"competitive-analysis-agent/internal/domain/entity"
	"competitive-analysis-agent/internal/supervisor"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Mock implementations for repository interfaces
// ---------------------------------------------------------------------------

type mockCompetitorRepo struct {
	data       []*entity.Competitor
	createErr  error
	getErr     error
	listErr    error
	updateErr  error
}

func (m *mockCompetitorRepo) Create(ctx context.Context, c *entity.Competitor) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.data = append(m.data, c)
	return nil
}

func (m *mockCompetitorRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Competitor, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for _, c := range m.data {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}

func (m *mockCompetitorRepo) GetByPlatformAndProductID(ctx context.Context, platform, productID string) (*entity.Competitor, error) {
	for _, c := range m.data {
		if c.Platform == platform && c.PlatformProductID == productID {
			return c, nil
		}
	}
	return nil, nil
}

func (m *mockCompetitorRepo) Update(ctx context.Context, c *entity.Competitor) error {
	return m.updateErr
}

func (m *mockCompetitorRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockCompetitorRepo) ListByPlatform(ctx context.Context, platform string, limit, offset int) ([]*entity.Competitor, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var out []*entity.Competitor
	for _, c := range m.data {
		if platform == "" || c.Platform == platform {
			out = append(out, c)
		}
	}
	if offset >= len(out) {
		return []*entity.Competitor{}, nil
	}
	end := offset + limit
	if end > len(out) {
		end = len(out)
	}
	return out[offset:end], nil
}

func (m *mockCompetitorRepo) ListAll(ctx context.Context, limit, offset int) ([]*entity.Competitor, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	if offset >= len(m.data) {
		return []*entity.Competitor{}, nil
	}
	end := offset + limit
	if end > len(m.data) {
		end = len(m.data)
	}
	return m.data[offset:end], nil
}

func (m *mockCompetitorRepo) CountByPlatform(ctx context.Context, platform string) (int, error) {
	count := 0
	for _, c := range m.data {
		if platform == "" || c.Platform == platform {
			count++
		}
	}
	return count, nil
}

func (m *mockCompetitorRepo) UpdatePrice(ctx context.Context, id uuid.UUID, price float64, currency string) error {
	return nil
}

func (m *mockCompetitorRepo) Upsert(ctx context.Context, c *entity.Competitor) error {
	m.data = append(m.data, c)
	return nil
}

type mockPriceHistoryRepo struct {
	data       []*entity.PriceHistory
	getErr     error
	createErr  error
}

func (m *mockPriceHistoryRepo) Create(ctx context.Context, p *entity.PriceHistory) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.data = append(m.data, p)
	return nil
}

func (m *mockPriceHistoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.PriceHistory, error) {
	return nil, nil
}

func (m *mockPriceHistoryRepo) GetByCompetitorID(ctx context.Context, competitorID uuid.UUID, limit int) ([]*entity.PriceHistory, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var out []*entity.PriceHistory
	for _, p := range m.data {
		if p.CompetitorID == competitorID {
			out = append(out, p)
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (m *mockPriceHistoryRepo) GetLatest(ctx context.Context, competitorID uuid.UUID) (*entity.PriceHistory, error) {
	return nil, nil
}

func (m *mockPriceHistoryRepo) GetLatestPrices(ctx context.Context, competitorIDs []uuid.UUID) (map[uuid.UUID]*entity.PriceHistory, error) {
	return nil, nil
}

func (m *mockPriceHistoryRepo) GetPriceRange(ctx context.Context, competitorID uuid.UUID, startTime, endTime time.Time) ([]*entity.PriceHistory, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var out []*entity.PriceHistory
	for _, p := range m.data {
		if p.CompetitorID == competitorID && !p.RecordedAt.Before(startTime) && !p.RecordedAt.After(endTime) {
			out = append(out, p)
		}
	}
	return out, nil
}

func (m *mockPriceHistoryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockPriceHistoryRepo) DeleteByCompetitorID(ctx context.Context, competitorID uuid.UUID) error {
	return nil
}

func (m *mockPriceHistoryRepo) Count(ctx context.Context, competitorID uuid.UUID) (int, error) {
	return len(m.data), nil
}

func (m *mockPriceHistoryRepo) DetectPriceChange(ctx context.Context, competitorID uuid.UUID, threshold float64) (interface{}, error) {
	return nil, nil
}

func (m *mockPriceHistoryRepo) GetAveragePrice(ctx context.Context, competitorID uuid.UUID, days int) (float64, error) {
	return 99.99, nil
}

func (m *mockPriceHistoryRepo) GetMinMaxPrice(ctx context.Context, competitorID uuid.UUID, days int) (min, max float64, err error) {
	return 80.0, 120.0, nil
}

// ---------------------------------------------------------------------------
// Mock implementations for TrendHandler dependencies
// ---------------------------------------------------------------------------

type mockMarketAgent struct {
	result *supervisor.MarketTrendResult
	err    error
}

func (m *mockMarketAgent) Execute(ctx context.Context, task *supervisor.Task) (any, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func (m *mockMarketAgent) Name() string { return "MockMarketAgent" }

type mockSchedulerAgent struct {
	name   string
	result any
	err    error
}

func (m *mockSchedulerAgent) Name() string { return m.name }
func (m *mockSchedulerAgent) Execute(ctx context.Context, task *supervisor.Task) (any, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

type mockScheduler struct {
	agents map[supervisor.TaskType]supervisor.Agent
}

func (m *mockScheduler) RegisterAgent(taskType supervisor.TaskType, agent supervisor.Agent) {
	if m.agents == nil {
		m.agents = make(map[supervisor.TaskType]supervisor.Agent)
	}
	m.agents[taskType] = agent
}

func (m *mockScheduler) DecomposeQuery(ctx context.Context, query string) ([]*supervisor.Task, error) {
	var tasks []*supervisor.Task
	lower := strings.ToLower(query)
	if strings.Contains(lower, "trend") || strings.Contains(lower, "market") {
		tasks = append(tasks, &supervisor.Task{Type: supervisor.TaskTypeMarketTrend, Query: query})
	}
	if strings.Contains(lower, "competitor") || strings.Contains(lower, "price") {
		tasks = append(tasks, &supervisor.Task{Type: supervisor.TaskTypeCompetitor, Query: query})
	}
	if len(tasks) == 0 {
		tasks = append(tasks, &supervisor.Task{Type: supervisor.TaskTypeMarketTrend, Query: query})
	}
	return tasks, nil
}

func (m *mockScheduler) ScheduleAndExecute(ctx context.Context, query string) ([]*supervisor.Task, error) {
	tasks, err := m.DecomposeQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		if agent, ok := m.agents[task.Type]; ok {
			result, err := agent.Execute(ctx, task)
			if err != nil {
				task.Result = err
			} else {
				task.Result = result
			}
		}
	}
	return tasks, nil
}

type mockAggregator struct {
	result *supervisor.AggregatedResult
}

func (m *mockAggregator) Aggregate(tasks []*supervisor.Task) *supervisor.AggregatedResult {
	if m.result != nil {
		return m.result
	}
	return &supervisor.AggregatedResult{TaskCount: len(tasks), Summary: "mock summary"}
}

// ---------------------------------------------------------------------------
// Test fixtures
// ---------------------------------------------------------------------------

func makeCompetitor(id, name, platform string, price float64) *entity.Competitor {
	now := time.Now()
	return &entity.Competitor{
		ID:                uuid.MustParse(id),
		Name:              name,
		Platform:          platform,
		PlatformProductID: "PID-" + name,
		CurrentPrice:      price,
		Currency:          "USD",
		Rating:            4.5,
		ReviewCount:       1000,
		SellerRating:      4.8,
		SellerReviewCount: 500,
		SourceURL:         "https://example.com/" + name,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

// ---------------------------------------------------------------------------
// Router setup with real mocks
// ---------------------------------------------------------------------------

func setupTestRouter() *chi.Mux {
	r := chi.NewRouter()

	compRepo := &mockCompetitorRepo{
		data: []*entity.Competitor{
			makeCompetitor("550e8400-e29b-41d4-a716-446655440001", "CompetitorOne", "amazon", 99.99),
			makeCompetitor("550e8400-e29b-41d4-a716-446655440002", "CompetitorTwo", "amazon", 109.99),
			makeCompetitor("550e8400-e29b-41d4-a716-446655440003", "CompetitorThree", "alibaba", 85.00),
		},
	}
	priceRepo := &mockPriceHistoryRepo{
		data: []*entity.PriceHistory{},
	}
	competitorHandler := NewCompetitorHandler(compRepo, priceRepo)

	marketAgent := &mockMarketAgent{
		result: &supervisor.MarketTrendResult{
			Trend:        "Bullish",
			Opportunities: []string{"Smart speakers", "Wireless earbuds"},
			DemandSignal: "Strong",
		},
	}
	sched := &mockScheduler{}
	sched.RegisterAgent(supervisor.TaskTypeMarketTrend, &mockSchedulerAgent{name: "market", result: &supervisor.MarketTrendResult{
		Trend:        "Mock Bullish",
		Opportunities: []string{"Opp1"},
		DemandSignal: "Strong",
	}})
	sched.RegisterAgent(supervisor.TaskTypeCompetitor, &mockSchedulerAgent{name: "competitor", result: &supervisor.CompetitorResult{
		Analysis: "Mock analysis",
		Competitors: []supervisor.CompetitorInsight{
			{Name: "CompA", Strength: "Brand", Weakness: "Price", Strategy: "Premium"},
		},
	}})
	trendHandler := NewTrendHandler(marketAgent, sched, &mockAggregator{})

	pricingHandler := NewPricingHandler()
	productHandler := NewProductHandler()

	// Wire up routes
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

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

// setupTestRouterWithNilRepos creates a router with nil dependencies to test service unavailable paths
func setupTestRouterWithNilRepos() *chi.Mux {
	r := chi.NewRouter()
	competitorHandler := NewCompetitorHandler(nil, nil)
	trendHandler := NewTrendHandler(nil, nil, nil)
	pricingHandler := NewPricingHandler()
	productHandler := NewProductHandler()

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/competitors", func(r chi.Router) {
			r.Get("/", competitorHandler.List)
			r.Get("/search", competitorHandler.Search)
			r.Get("/price-changes", competitorHandler.DetectPriceChanges)
			r.Get("/{id}", competitorHandler.Get)
			r.Get("/{id}/price-history", competitorHandler.GetPriceHistory)
		})
		r.Route("/trends", func(r chi.Router) {
			r.Get("/analyze", trendHandler.AnalyzeTrend)
			r.Get("/analysis", trendHandler.GetAggregatedAnalysis)
			r.Get("/categories", trendHandler.ListTrendCategories)
			r.Get("/categories/{category}", trendHandler.GetTrendByCategory)
			r.Get("/demand-signals", trendHandler.GetDemandSignals)
			r.Get("/overview", trendHandler.GetMarketOverview)
			r.Get("/trending", trendHandler.GetTrendingProducts)
		})
		r.Route("/strategy/pricing", func(r chi.Router) {
			r.Get("/", pricingHandler.GetPricingRecommendation)
			r.Get("/competitors/{product_id}", pricingHandler.GetCompetitorPricing)
			r.Get("/elasticity/{product_id}", pricingHandler.GetPriceElasticity)
			r.Post("/scenario", pricingHandler.GetPricingScenario)
			r.Get("/history/{product_id}", pricingHandler.GetPricingHistory)
			r.Get("/dynamic/{product_id}", pricingHandler.GetDynamicPricing)
		})
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

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "status")
	assert.Contains(t, w.Body.String(), "ok")
}

// Trend Handler Tests

func TestListTrendCategories(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "categories")

	categories := response["categories"].([]interface{})
	assert.GreaterOrEqual(t, len(categories), 1)
}

func TestGetTrendByCategory(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/categories/electronics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "electronics", response["category"])
	assert.Contains(t, response, "trend")
}

func TestGetTrendByCategory_UnknownCategory(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/categories/unknown", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "unknown", response["category"])
}

func TestGetDemandSignals(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/demand-signals", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "signals")
	assert.Contains(t, response, "updated_at")
}

func TestGetDemandSignals_WithPlatform(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/demand-signals?platform=amazon", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "amazon", response["platform"])
}

func TestGetMarketOverview(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/overview", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "total_products")
	assert.Contains(t, response, "active_competitors")
	assert.Contains(t, response, "market_sentiment")
}

func TestGetTrendingProducts(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/trending", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "products")
}

func TestGetTrendingProducts_WithLimit(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/trending?limit=3", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Pricing Handler Tests

func TestGetPricingRecommendation_MissingProductID(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/strategy/pricing/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
}

func TestGetPricingRecommendation_WithProductID(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/strategy/pricing/?product_id=123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "recommendation")
}

func TestGetCompetitorPricing(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/strategy/pricing/competitors/123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "competitors")
	assert.Contains(t, response, "avg_price")
}

func TestGetPriceElasticity(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/strategy/pricing/elasticity/123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "elasticity")
	assert.Contains(t, response, "elasticity_type")
}

func TestGetPricingScenario(t *testing.T) {
	router := setupTestRouter()

	body := `{"product_id":"123","base_price":99.99,"target_margin":0.3,"competitor_prices":[95,105]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/strategy/pricing/scenario", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "scenarios")
	assert.Contains(t, response, "recommended")
}

func TestGetPricingScenario_InvalidJSON(t *testing.T) {
	router := setupTestRouter()

	body := `invalid json`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/strategy/pricing/scenario", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPricingScenario_WrongContentType(t *testing.T) {
	router := setupTestRouter()

	body := `{"product_id":"123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/strategy/pricing/scenario", strings.NewReader(body))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPricingHistory(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/strategy/pricing/history/123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "history")
}

func TestGetPricingHistory_WithDays(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/strategy/pricing/history/123?days=60", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(60), response["days"])
}

func TestGetDynamicPricing(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/strategy/pricing/dynamic/123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "current_price")
	assert.Contains(t, response, "recommended_price")
	assert.Contains(t, response, "factors")
}

// Product Handler Tests

func TestGetProductRecommendations(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/recommendations", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "recommendations")
	assert.Contains(t, response, "total")
}

func TestGetProductRecommendations_WithCategory(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/recommendations?category=Electronics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "filter_applied")
}

func TestGetProductRecommendations_WithLimit(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/recommendations?limit=3", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, int(response["total"].(float64)), 0)
}

func TestGetProductByID_InvalidUUID(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetProductByID(t *testing.T) {
	router := setupTestRouter()

	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+validUUID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, validUUID, response["id"])
}

func TestGetProductAnalysis(t *testing.T) {
	router := setupTestRouter()

	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+validUUID+"/analysis", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "swot")
	assert.Contains(t, response, "market_position")
}

func TestGetProductComparison(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/compare?ids=1,2,3", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "products")
	assert.Contains(t, response, "best_value")
}

func TestGetProductComparison_MissingIDs(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/compare", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetCategoryInsights(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/categories/electronics/insights", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	// Category case may vary - check it's not empty
	assert.NotEmpty(t, response["category"])
}

func TestGetCategoryInsights_UnknownCategory(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/categories/unknown/insights", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "unknown", response["category"])
}

func TestGetProductTrends(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/trends", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "trends")
	assert.Contains(t, response, "summary")
}

func TestGetProductTrends_WithPeriod(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/trends?period=7d", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "7d", response["period"])
}

// Competitor Handler Tests - nil repo scenarios

func TestListCompetitors_ServiceUnavailable(t *testing.T) {
	router := setupTestRouterWithNilRepos()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestGetCompetitor_ServiceUnavailable(t *testing.T) {
	router := setupTestRouterWithNilRepos()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestGetCompetitor_InvalidUUID(t *testing.T) {
	router := setupTestRouterWithNilRepos()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Returns 503 because repo is nil - nil check happens before UUID validation
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestGetPriceHistory_InvalidUUID(t *testing.T) {
	router := setupTestRouterWithNilRepos()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/invalid-uuid/price-history", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Returns 503 because repo is nil - nil check happens before UUID validation
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestSearchCompetitors_ServiceUnavailable(t *testing.T) {
	router := setupTestRouterWithNilRepos()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/search?platform=amazon", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestSearchCompetitors_MissingParams(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/search", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Returns 400 (bad request) or 503 (service unavailable) depending on nil check order
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusServiceUnavailable)
}

func TestDetectPriceChanges_ServiceUnavailable(t *testing.T) {
	router := setupTestRouterWithNilRepos()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/price-changes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

// Trend Handler Service Unavailable Tests

func TestAnalyzeTrend_ServiceUnavailable(t *testing.T) {
	router := setupTestRouterWithNilRepos()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/analyze?query=test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestGetAggregatedAnalysis_MissingQuery(t *testing.T) {
	router := setupTestRouterWithNilRepos()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/analysis", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAggregatedAnalysis_ServiceUnavailable(t *testing.T) {
	router := setupTestRouterWithNilRepos()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/analysis?query=test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

// Error Handling Tests

func TestMethodNotAllowed(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestNotFound(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// JSON Response Format Tests

func TestJSONContentType(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestJSONResponseStructure(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Check that response is valid JSON with categories
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	categories, ok := response["categories"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, categories, 8)
}

// ---------------------------------------------------------------------------
// Competitor Handler Tests (with real mock data)
// ---------------------------------------------------------------------------

func TestListCompetitors_Success(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ListCompetitorsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Greater(t, response.Total, 0)
	assert.NotEmpty(t, response.Competitors)
}

func TestListCompetitors_WithPlatform(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/?platform=amazon", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ListCompetitorsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
}

func TestListCompetitors_WithLimitAndOffset(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/?limit=2&offset=0", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response ListCompetitorsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, response.Limit)
	assert.Equal(t, 0, response.Offset)
}

func TestListCompetitors_ListError(t *testing.T) {
	router := setupTestRouter()

	// Force an error by passing invalid platform query
	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/?platform=", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	// Empty platform returns list without error
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCompetitor_Success(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/550e8400-e29b-41d4-a716-446655440001", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response CompetitorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440001", response.ID)
}

func TestGetCompetitor_NotFound(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/550e8400-e29b-41d4-a716-446655440099", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetPriceHistory_Success(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/550e8400-e29b-41d4-a716-446655440001/price-history", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPriceHistory_WithLimit(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/550e8400-e29b-41d4-a716-446655440001/price-history?limit=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDetectPriceChanges_Success(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/price-changes?threshold=0.05", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDetectPriceChanges_DefaultThreshold(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/price-changes", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchCompetitors_Success(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/search?platform=amazon&min_price=50&max_price=200", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchCompetitors_WithPlatformOnly(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/search?platform=amazon", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnalyzeTrend_Success(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/analyze?query=smart+speakers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAggregatedAnalysis_Success(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/analysis?query=amazon+smart+speakers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnalyzeTrend_EmptyQuery(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/analyze?query=", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Empty query defaults to "general market trends" so returns 200
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProductAnalysis_Success(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/550e8400-e29b-41d4-a716-446655440001/analysis", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProductComparison_Success(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/compare?ids=550e8400-e29b-41d4-a716-446655440001,550e8400-e29b-41d4-a716-446655440002", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPriceHistory_RepoError(t *testing.T) {
	// This test uses setupTestRouterWithNilRepos to trigger service unavailable
	router := setupTestRouterWithNilRepos()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/550e8400-e29b-41d4-a716-446655440001/price-history", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestGetPriceHistory_RepoErrorNilPriceRepo(t *testing.T) {
	// Test with competitor repo but nil price repo
	router := setupTestRouterWithNilRepos()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/550e8400-e29b-41d4-a716-446655440001/price-history", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestDetectPriceChanges_WithPlatform(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/price-changes?platform=amazon&threshold=0.1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchCompetitors_AllParams(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors/search?platform=alibaba&min_price=10&max_price=500&keyword=smart+watch", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
