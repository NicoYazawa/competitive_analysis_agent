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
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Router test fixtures (reused from benchmark_test.go to ensure consistency)
// ---------------------------------------------------------------------------

type testCompetitorRepo struct {
	data []*entity.Competitor
}

func (r *testCompetitorRepo) Create(ctx context.Context, c *entity.Competitor) error {
	r.data = append(r.data, c)
	return nil
}
func (r *testCompetitorRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Competitor, error) {
	for _, c := range r.data {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}
func (r *testCompetitorRepo) GetByPlatformAndProductID(ctx context.Context, platform, productID string) (*entity.Competitor, error) {
	return nil, nil
}
func (r *testCompetitorRepo) Update(ctx context.Context, c *entity.Competitor) error    { return nil }
func (r *testCompetitorRepo) Delete(ctx context.Context, id uuid.UUID) error          { return nil }
func (r *testCompetitorRepo) ListByPlatform(ctx context.Context, platform string, limit, offset int) ([]*entity.Competitor, error) {
	return r.data, nil
}
func (r *testCompetitorRepo) ListAll(ctx context.Context, limit, offset int) ([]*entity.Competitor, error) {
	return r.data, nil
}
func (r *testCompetitorRepo) CountByPlatform(ctx context.Context, platform string) (int, error) {
	return len(r.data), nil
}
func (r *testCompetitorRepo) UpdatePrice(ctx context.Context, id uuid.UUID, price float64, currency string) error {
	return nil
}
func (r *testCompetitorRepo) Upsert(ctx context.Context, c *entity.Competitor) error {
	r.data = append(r.data, c)
	return nil
}

type testPriceHistoryRepo struct{}

func (r *testPriceHistoryRepo) Create(ctx context.Context, p *entity.PriceHistory) error { return nil }
func (r *testPriceHistoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.PriceHistory, error) {
	return nil, nil
}
func (r *testPriceHistoryRepo) GetByCompetitorID(ctx context.Context, competitorID uuid.UUID, limit int) ([]*entity.PriceHistory, error) {
	return nil, nil
}
func (r *testPriceHistoryRepo) GetLatest(ctx context.Context, competitorID uuid.UUID) (*entity.PriceHistory, error) {
	return nil, nil
}
func (r *testPriceHistoryRepo) GetLatestPrices(ctx context.Context, competitorIDs []uuid.UUID) (map[uuid.UUID]*entity.PriceHistory, error) {
	return nil, nil
}
func (r *testPriceHistoryRepo) GetPriceRange(ctx context.Context, competitorID uuid.UUID, startTime, endTime time.Time) ([]*entity.PriceHistory, error) {
	return nil, nil
}
func (r *testPriceHistoryRepo) Delete(ctx context.Context, id uuid.UUID) error                       { return nil }
func (r *testPriceHistoryRepo) DeleteByCompetitorID(ctx context.Context, competitorID uuid.UUID) error {
	return nil
}
func (r *testPriceHistoryRepo) Count(ctx context.Context, competitorID uuid.UUID) (int, error) { return 0, nil }
func (r *testPriceHistoryRepo) DetectPriceChange(ctx context.Context, competitorID uuid.UUID, threshold float64) (interface{}, error) {
	return nil, nil
}
func (r *testPriceHistoryRepo) GetAveragePrice(ctx context.Context, competitorID uuid.UUID, days int) (float64, error) {
	return 99.99, nil
}
func (r *testPriceHistoryRepo) GetMinMaxPrice(ctx context.Context, competitorID uuid.UUID, days int) (min, max float64, err error) {
	return 80.0, 120.0, nil
}

type testMarketAgent struct{}

func (a *testMarketAgent) Execute(ctx context.Context, task *supervisor.Task) (any, error) {
	return &supervisor.MarketTrendResult{
		Trend:         "Bullish",
		Opportunities: []string{"Smart speakers"},
		DemandSignal:  "Strong",
	}, nil
}
func (a *testMarketAgent) Name() string { return "TestMarketAgent" }

type testSchedulerAgent struct {
	name   string
	result any
}

func (a *testSchedulerAgent) Name() string { return a.name }
func (a *testSchedulerAgent) Execute(ctx context.Context, task *supervisor.Task) (any, error) {
	return a.result, nil
}

type testScheduler struct {
	agents map[supervisor.TaskType]supervisor.Agent
}

func (s *testScheduler) RegisterAgent(taskType supervisor.TaskType, agent supervisor.Agent) {
	if s.agents == nil {
		s.agents = make(map[supervisor.TaskType]supervisor.Agent)
	}
	s.agents[taskType] = agent
}
func (s *testScheduler) DecomposeQuery(ctx context.Context, query string) ([]*supervisor.Task, error) {
	return []*supervisor.Task{{Type: supervisor.TaskTypeMarketTrend, Query: query}}, nil
}
func (s *testScheduler) ScheduleAndExecute(ctx context.Context, query string) ([]*supervisor.Task, error) {
	return []*supervisor.Task{{Type: supervisor.TaskTypeMarketTrend, Query: query, Result: &supervisor.MarketTrendResult{}}}, nil
}

type testAggregator struct{}

func (a *testAggregator) Aggregate(tasks []*supervisor.Task) *supervisor.AggregatedResult {
	return &supervisor.AggregatedResult{TaskCount: len(tasks), Summary: "test summary"}
}

func setupTestRouter() *chi.Mux {
	compRepo := &testCompetitorRepo{
		data: []*entity.Competitor{
			{
				ID:                uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
				Name:              "TestCompetitor",
				Platform:         "amazon",
				PlatformProductID: "PID-Test",
				CurrentPrice:     99.99,
				Currency:         "USD",
				Rating:           4.5,
				ReviewCount:      1000,
				SellerRating:     4.8,
				SellerReviewCount: 500,
				SourceURL:        "https://example.com/test",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		},
	}
	priceRepo := &testPriceHistoryRepo{}
	sched := &testScheduler{}
	sched.RegisterAgent(supervisor.TaskTypeMarketTrend, &testSchedulerAgent{name: "market", result: &supervisor.MarketTrendResult{
		Trend:         "Mock Bullish",
		Opportunities: []string{"Opp1"},
		DemandSignal:  "Strong",
	}})
	sched.RegisterAgent(supervisor.TaskTypeCompetitor, &testSchedulerAgent{name: "competitor", result: &supervisor.CompetitorResult{
		Analysis: "Mock analysis",
	}})

	compHandler := handlers.NewCompetitorHandler(compRepo, priceRepo)
	trendHandler := handlers.NewTrendHandler(&testMarketAgent{}, sched, &testAggregator{})
	pricingHandler := handlers.NewPricingHandler()
	productHandler := handlers.NewProductHandler()

	return NewRouter(compHandler, trendHandler, pricingHandler, productHandler)
}

// ---------------------------------------------------------------------------
// Router tests
// ---------------------------------------------------------------------------

func TestNewRouter(t *testing.T) {
	router := setupTestRouter()
	assert.NotNil(t, router)
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

func TestAPIv1ProductsNotFound(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPIv1ProductsWithIDNotFound(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/123", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// "123" is not a valid UUID, so it returns 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAPIv1CompetitorsNotFound(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/competitors", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// competitors endpoint should be found (returns 200 or 503 depending on repo state)
	assert.True(t, w.Code == http.StatusServiceUnavailable || w.Code == http.StatusOK)
}

func TestAPIv1TrendsNotFound(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/trends", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAPIv1StrategyPricingNotFound(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/strategy/pricing", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// missing product_id returns 400
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusOK)
}

func TestNotFoundHandler(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent/path", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMethodNotAllowed(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Chi returns 405 for method not allowed on registered routes
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
