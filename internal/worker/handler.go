package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"competitive-analysis-agent/internal/scraper"

	"github.com/hibiken/asynq"
)

// ProductScraper defines the interface for scraping competitor product data.
type ProductScraper interface {
	ScrapePrices(ctx context.Context, productID string) ([]float64, error)
	ScrapeCompetitors(ctx context.Context, query string) ([]string, error)
}

// LLMClient defines the interface for LLM-based analysis.
type LLMClient interface {
	AnalyzeTrends(ctx context.Context, query string) (map[string]interface{}, error)
	CheckSupplyChain(ctx context.Context, productID string) ([]string, error)
}

// Handler processes background tasks for competitive analysis.
type Handler struct {
	logger  *slog.Logger
	scraper ProductScraper
	llm     LLMClient
}

// NewHandler creates a Handler with real scraper and LLM implementations.
// Both dependencies are optional; if nil the handler uses fallback behavior
// that is suitable only for development/testing.
func NewHandler(logger *slog.Logger, scraper ProductScraper, llm LLMClient) *Handler {
	return &Handler{
		logger:  logger,
		scraper: scraper,
		llm:     llm,
	}
}

// HandlePriceCheck processes a price check task.
func (h *Handler) HandlePriceCheck(ctx context.Context, payload *TaskPayload) error {
	h.logger.Info("Processing price check task",
		slog.String("task_id", payload.TaskID),
		slog.String("product_id", payload.ProductID))

	var prices []float64
	var err error

	if h.scraper != nil {
		prices, err = h.scraper.ScrapePrices(ctx, payload.ProductID)
	} else {
		prices, err = h.devScrapePrices(ctx, payload.ProductID)
	}
	if err != nil {
		return fmt.Errorf("scrape prices failed: %w", err)
	}

	h.logger.Info("Price check completed",
		slog.String("task_id", payload.TaskID),
		slog.Int("prices_count", len(prices)))

	return nil
}

// HandleCompetitorSync processes a competitor data sync task.
func (h *Handler) HandleCompetitorSync(ctx context.Context, payload *TaskPayload) error {
	h.logger.Info("Processing competitor sync task",
		slog.String("task_id", payload.TaskID),
		slog.String("query", payload.Query))

	var competitors []string
	var err error

	if h.scraper != nil {
		competitors, err = h.scraper.ScrapeCompetitors(ctx, payload.Query)
	} else {
		competitors, err = h.devSyncCompetitorData(ctx, payload.Query)
	}
	if err != nil {
		return fmt.Errorf("sync competitors failed: %w", err)
	}

	h.logger.Info("Competitor sync completed",
		slog.String("task_id", payload.TaskID),
		slog.Int("competitors_count", len(competitors)))

	return nil
}

// HandleTrendAnalysis processes a market trend analysis task.
func (h *Handler) HandleTrendAnalysis(ctx context.Context, payload *TaskPayload) error {
	h.logger.Info("Processing trend analysis task",
		slog.String("task_id", payload.TaskID),
		slog.String("query", payload.Query))

	if h.llm == nil {
		return fmt.Errorf("trend analysis requires LLM client")
	}

	trends, err := h.llm.AnalyzeTrends(ctx, payload.Query)
	if err != nil {
		return fmt.Errorf("analyze trends failed: %w", err)
	}

	h.logger.Info("Trend analysis completed",
		slog.String("task_id", payload.TaskID),
		slog.Any("trends", trends))

	return nil
}

// HandleSupplyAlert processes a supply chain alert task.
func (h *Handler) HandleSupplyAlert(ctx context.Context, payload *TaskPayload) error {
	h.logger.Info("Processing supply alert task",
		slog.String("task_id", payload.TaskID),
		slog.String("product_id", payload.ProductID))

	if h.llm == nil {
		return fmt.Errorf("supply alert requires LLM client")
	}

	alerts, err := h.llm.CheckSupplyChain(ctx, payload.ProductID)
	if err != nil {
		return fmt.Errorf("check supply chain failed: %w", err)
	}

	h.logger.Info("Supply alert completed",
		slog.String("task_id", payload.TaskID),
		slog.Int("alerts_count", len(alerts)))

	return nil
}

// ProcessTask is the entry point for all task types.
func (h *Handler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	var payload TaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("parse payload failed: %w", err)
	}

	h.logger.Info("Processing task",
		slog.String("type", string(payload.Type)),
		slog.String("task_id", payload.TaskID))

	switch payload.Type {
	case TaskTypePriceCheck:
		return h.HandlePriceCheck(ctx, &payload)
	case TaskTypeCompetitorSync:
		return h.HandleCompetitorSync(ctx, &payload)
	case TaskTypeTrendAnalysis:
		return h.HandleTrendAnalysis(ctx, &payload)
	case TaskTypeSupplyAlert:
		return h.HandleSupplyAlert(ctx, &payload)
	default:
		return fmt.Errorf("unknown task type: %s", payload.Type)
	}
}

// devScrapePrices is a development fallback that returns hardcoded price data.
// It is only used when no ProductScraper is configured.
func (h *Handler) devScrapePrices(ctx context.Context, productID string) ([]float64, error) {
	h.logger.Warn("using dev fallback for ScrapePrices — replace with real scraper for production")
	time.Sleep(100 * time.Millisecond)
	return []float64{99.99, 109.99, 95.50, 105.00}, nil
}

// devSyncCompetitorData is a development fallback that returns hardcoded competitor names.
// It is only used when no ProductScraper is configured.
func (h *Handler) devSyncCompetitorData(ctx context.Context, query string) ([]string, error) {
	h.logger.Warn("using dev fallback for ScrapeCompetitors — replace with real scraper for production")
	time.Sleep(100 * time.Millisecond)
	return []string{"CompetitorA", "CompetitorB", "CompetitorC"}, nil
}

// ---------------------------------------------------------------------------
// Real implementations — used when dependencies are injected
// ---------------------------------------------------------------------------

// ScrapeCompetitorPricesViaScraper wraps a CollyScraper to implement ProductScraper.
// It extracts prices from Amazon product pages derived from the product ID.
func ScrapeCompetitorPricesViaScraper(s *scraper.CollyScraper, productID string, urls []string) ([]float64, error) {
	results := s.ScrapeMultiple(context.Background(), urls)
	var prices []float64
	for _, r := range results {
		if r.Error != nil {
			continue
		}
		if price, err := scraper.ParsePrice(r.Price); err == nil {
			prices = append(prices, price)
		}
	}
	return prices, nil
}
