package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"competitive-analysis-agent/internal/api/handlers"
	"competitive-analysis-agent/internal/domain/entity"
	"competitive-analysis-agent/internal/supervisor"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// Benchmark test fixtures and helpers
// ---------------------------------------------------------------------------

func makeBenchmarkCompetitor(id, name, platform string, price float64) *entity.Competitor {
	now := time.Now()
	return &entity.Competitor{
		ID:                uuid.MustParse(id),
		Name:              name,
		Platform:         platform,
		PlatformProductID: "PID-" + name,
		CurrentPrice:      price,
		Currency:         "USD",
		Rating:           4.5,
		ReviewCount:      1000,
		SellerRating:     4.8,
		SellerReviewCount: 500,
		SourceURL:        "https://example.com/" + name,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

type benchmarkCompetitorRepo struct {
	data []*entity.Competitor
}

func (r *benchmarkCompetitorRepo) Create(ctx context.Context, c *entity.Competitor) error {
	r.data = append(r.data, c)
	return nil
}
func (r *benchmarkCompetitorRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Competitor, error) {
	for _, c := range r.data {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}
func (r *benchmarkCompetitorRepo) GetByPlatformAndProductID(ctx context.Context, platform, productID string) (*entity.Competitor, error) {
	return nil, nil
}
func (r *benchmarkCompetitorRepo) Update(ctx context.Context, c *entity.Competitor) error    { return nil }
func (r *benchmarkCompetitorRepo) Delete(ctx context.Context, id uuid.UUID) error          { return nil }
func (r *benchmarkCompetitorRepo) ListByPlatform(ctx context.Context, platform string, limit, offset int) ([]*entity.Competitor, error) {
	return r.data, nil
}
func (r *benchmarkCompetitorRepo) ListAll(ctx context.Context, limit, offset int) ([]*entity.Competitor, error) {
	return r.data, nil
}
func (r *benchmarkCompetitorRepo) CountByPlatform(ctx context.Context, platform string) (int, error) {
	return len(r.data), nil
}
func (r *benchmarkCompetitorRepo) UpdatePrice(ctx context.Context, id uuid.UUID, price float64, currency string) error {
	return nil
}
func (r *benchmarkCompetitorRepo) Upsert(ctx context.Context, c *entity.Competitor) error {
	r.data = append(r.data, c)
	return nil
}

type benchmarkPriceHistoryRepo struct{}

func (r *benchmarkPriceHistoryRepo) Create(ctx context.Context, p *entity.PriceHistory) error { return nil }
func (r *benchmarkPriceHistoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.PriceHistory, error) {
	return nil, nil
}
func (r *benchmarkPriceHistoryRepo) GetByCompetitorID(ctx context.Context, competitorID uuid.UUID, limit int) ([]*entity.PriceHistory, error) {
	return nil, nil
}
func (r *benchmarkPriceHistoryRepo) GetLatest(ctx context.Context, competitorID uuid.UUID) (*entity.PriceHistory, error) {
	return nil, nil
}
func (r *benchmarkPriceHistoryRepo) GetLatestPrices(ctx context.Context, competitorIDs []uuid.UUID) (map[uuid.UUID]*entity.PriceHistory, error) {
	return nil, nil
}
func (r *benchmarkPriceHistoryRepo) GetPriceRange(ctx context.Context, competitorID uuid.UUID, startTime, endTime time.Time) ([]*entity.PriceHistory, error) {
	return nil, nil
}
func (r *benchmarkPriceHistoryRepo) Delete(ctx context.Context, id uuid.UUID) error                       { return nil }
func (r *benchmarkPriceHistoryRepo) DeleteByCompetitorID(ctx context.Context, competitorID uuid.UUID) error { return nil }
func (r *benchmarkPriceHistoryRepo) Count(ctx context.Context, competitorID uuid.UUID) (int, error)        { return 0, nil }
func (r *benchmarkPriceHistoryRepo) DetectPriceChange(ctx context.Context, competitorID uuid.UUID, threshold float64) (interface{}, error) {
	return nil, nil
}
func (r *benchmarkPriceHistoryRepo) GetAveragePrice(ctx context.Context, competitorID uuid.UUID, days int) (float64, error) {
	return 99.99, nil
}
func (r *benchmarkPriceHistoryRepo) GetMinMaxPrice(ctx context.Context, competitorID uuid.UUID, days int) (min, max float64, err error) {
	return 80.0, 120.0, nil
}

type benchmarkMarketAgent struct{}

func (a *benchmarkMarketAgent) Execute(ctx context.Context, task *supervisor.Task) (any, error) {
	return &supervisor.MarketTrendResult{
		Trend:         "Bullish",
		Opportunities: []string{"Smart speakers"},
		DemandSignal:  "Strong",
	}, nil
}
func (a *benchmarkMarketAgent) Name() string { return "BenchmarkMarketAgent" }

type benchmarkSchedulerAgent struct {
	name   string
	result any
}

func (a *benchmarkSchedulerAgent) Name() string { return a.name }
func (a *benchmarkSchedulerAgent) Execute(ctx context.Context, task *supervisor.Task) (any, error) {
	return a.result, nil
}

type benchmarkScheduler struct {
	agents map[supervisor.TaskType]supervisor.Agent
}

func (s *benchmarkScheduler) RegisterAgent(taskType supervisor.TaskType, agent supervisor.Agent) {
	if s.agents == nil {
		s.agents = make(map[supervisor.TaskType]supervisor.Agent)
	}
	s.agents[taskType] = agent
}
func (s *benchmarkScheduler) DecomposeQuery(ctx context.Context, query string) ([]*supervisor.Task, error) {
	return []*supervisor.Task{{Type: supervisor.TaskTypeMarketTrend, Query: query}}, nil
}
func (s *benchmarkScheduler) ScheduleAndExecute(ctx context.Context, query string) ([]*supervisor.Task, error) {
	return []*supervisor.Task{{Type: supervisor.TaskTypeMarketTrend, Query: query, Result: &supervisor.MarketTrendResult{}}}, nil
}

type benchmarkAggregator struct{}

func (a *benchmarkAggregator) Aggregate(tasks []*supervisor.Task) *supervisor.AggregatedResult {
	return &supervisor.AggregatedResult{TaskCount: len(tasks), Summary: "benchmark summary"}
}

func setupBenchmarkRouter() *chi.Mux {
	compRepo := &benchmarkCompetitorRepo{
		data: []*entity.Competitor{
			makeBenchmarkCompetitor("550e8400-e29b-41d4-a716-446655440001", "CompetitorOne", "amazon", 99.99),
			makeBenchmarkCompetitor("550e8400-e29b-41d4-a716-446655440002", "CompetitorTwo", "amazon", 109.99),
		},
	}
	priceRepo := &benchmarkPriceHistoryRepo{}
	sched := &benchmarkScheduler{}
	sched.RegisterAgent(supervisor.TaskTypeMarketTrend, &benchmarkSchedulerAgent{name: "market", result: &supervisor.MarketTrendResult{
		Trend:         "Mock Bullish",
		Opportunities: []string{"Opp1"},
		DemandSignal:  "Strong",
	}})
	sched.RegisterAgent(supervisor.TaskTypeCompetitor, &benchmarkSchedulerAgent{name: "competitor", result: &supervisor.CompetitorResult{
		Analysis: "Mock analysis",
	}})

	compHandler := handlers.NewCompetitorHandler(compRepo, priceRepo)
	trendHandler := handlers.NewTrendHandler(&benchmarkMarketAgent{}, sched, &benchmarkAggregator{})
	pricingHandler := handlers.NewPricingHandler()
	productHandler := handlers.NewProductHandler()

	return NewRouter(compHandler, trendHandler, pricingHandler, productHandler)
}

// ---------------------------------------------------------------------------
// Benchmark tests for P99 response time verification
// P99 Target: < 500ms for all endpoints
// ---------------------------------------------------------------------------

func BenchmarkHealthEndpoint(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkListTrendCategories(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/categories", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetTrendByCategory(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/categories/electronics", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetDemandSignals(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/demand-signals", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetMarketOverview(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/overview", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetTrendingProducts(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends/trending", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetPricingRecommendation(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/strategy/pricing/?product_id=123", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetCompetitorPricing(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/strategy/pricing/competitors/123", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetProductRecommendations(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/recommendations", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetProductComparison(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/compare?ids=1,2,3", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetCategoryInsights(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/categories/electronics/insights", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetProductTrends(b *testing.B) {
	router := setupBenchmarkRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/trends", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}
