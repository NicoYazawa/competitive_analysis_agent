package worker

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestNewHandler(t *testing.T) {
	logger := newTestLogger()
	handler := NewHandler(logger, nil, nil)

	assert.NotNil(t, handler)
	assert.Equal(t, logger, handler.logger)
}

func TestHandler_ProcessTask(t *testing.T) {
	logger := newTestLogger()
	handler := NewHandler(logger, nil, nil)

	tests := []struct {
		name        string
		taskType    TaskType
		productID   string
		query       string
		expectError bool
	}{
		{
			name:        "price check task",
			taskType:    TaskTypePriceCheck,
			productID:   "prod-001",
			expectError: false,
		},
		{
			name:        "competitor sync task",
			taskType:    TaskTypeCompetitorSync,
			query:       "smartphone",
			expectError: false,
		},
		{
			name:        "trend analysis task",
			taskType:    TaskTypeTrendAnalysis,
			query:       "market trends",
			expectError: true, // nil LLM returns error
		},
		{
			name:        "supply alert task",
			taskType:    TaskTypeSupplyAlert,
			productID:   "prod-002",
			expectError: true, // nil LLM returns error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := NewTaskPayload(tt.taskType, tt.productID, tt.query)
			payloadJSON, _ := payload.ToJSON()

			task := asynq.NewTask(string(tt.taskType), []byte(payloadJSON))
			err := handler.ProcessTask(context.Background(), task)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandler_ProcessTask_UnknownType(t *testing.T) {
	logger := newTestLogger()
	handler := NewHandler(logger, nil, nil)

	payload := NewTaskPayload(TaskType("unknown_type"), "", "")
	payloadJSON, _ := payload.ToJSON()

	task := asynq.NewTask("unknown_type", []byte(payloadJSON))
	err := handler.ProcessTask(context.Background(), task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown task type")
}

func TestHandler_ProcessTask_InvalidPayload(t *testing.T) {
	logger := newTestLogger()
	handler := NewHandler(logger, nil, nil)

	task := asynq.NewTask("price_check", []byte("invalid json"))
	err := handler.ProcessTask(context.Background(), task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse payload failed")
}

// ---------------------------------------------------------------------------
// Mock implementations for testing handler with real dependencies
// ---------------------------------------------------------------------------

type mockProductScraper struct {
	scrapePricesFunc      func(ctx context.Context, productID string) ([]float64, error)
	scrapeCompetitorsFunc func(ctx context.Context, query string) ([]string, error)
}

func (m *mockProductScraper) ScrapePrices(ctx context.Context, productID string) ([]float64, error) {
	if m.scrapePricesFunc != nil {
		return m.scrapePricesFunc(ctx, productID)
	}
	return nil, nil
}

func (m *mockProductScraper) ScrapeCompetitors(ctx context.Context, query string) ([]string, error) {
	if m.scrapeCompetitorsFunc != nil {
		return m.scrapeCompetitorsFunc(ctx, query)
	}
	return nil, nil
}

type mockLLMClient struct {
	analyzeTrendsFunc  func(ctx context.Context, query string) (map[string]interface{}, error)
	checkSupplyChainFunc func(ctx context.Context, productID string) ([]string, error)
}

func (m *mockLLMClient) AnalyzeTrends(ctx context.Context, query string) (map[string]interface{}, error) {
	if m.analyzeTrendsFunc != nil {
		return m.analyzeTrendsFunc(ctx, query)
	}
	return nil, nil
}

func (m *mockLLMClient) CheckSupplyChain(ctx context.Context, productID string) ([]string, error) {
	if m.checkSupplyChainFunc != nil {
		return m.checkSupplyChainFunc(ctx, productID)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// Tests for handler with mock dependencies
// ---------------------------------------------------------------------------

func TestHandler_HandlePriceCheck_WithScraper(t *testing.T) {
	logger := newTestLogger()

	t.Run("success", func(t *testing.T) {
		scraper := &mockProductScraper{
			scrapePricesFunc: func(ctx context.Context, productID string) ([]float64, error) {
				return []float64{99.99, 109.99}, nil
			},
		}
		handler := NewHandler(logger, scraper, nil)

		err := handler.HandlePriceCheck(context.Background(), NewTaskPayload(TaskTypePriceCheck, "prod-001", ""))
		assert.NoError(t, err)
	})

	t.Run("scraper error", func(t *testing.T) {
		scraper := &mockProductScraper{
			scrapePricesFunc: func(ctx context.Context, productID string) ([]float64, error) {
				return nil, fmt.Errorf("scraper error")
			},
		}
		handler := NewHandler(logger, scraper, nil)

		err := handler.HandlePriceCheck(context.Background(), NewTaskPayload(TaskTypePriceCheck, "prod-001", ""))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "scrape prices failed")
	})
}

func TestHandler_HandleCompetitorSync_WithScraper(t *testing.T) {
	logger := newTestLogger()

	t.Run("success", func(t *testing.T) {
		scraper := &mockProductScraper{
			scrapeCompetitorsFunc: func(ctx context.Context, query string) ([]string, error) {
				return []string{"CompetitorA", "CompetitorB"}, nil
			},
		}
		handler := NewHandler(logger, scraper, nil)

		err := handler.HandleCompetitorSync(context.Background(), NewTaskPayload(TaskTypeCompetitorSync, "", "smartphone"))
		assert.NoError(t, err)
	})

	t.Run("scraper error", func(t *testing.T) {
		scraper := &mockProductScraper{
			scrapeCompetitorsFunc: func(ctx context.Context, query string) ([]string, error) {
				return nil, fmt.Errorf("scraper error")
			},
		}
		handler := NewHandler(logger, scraper, nil)

		err := handler.HandleCompetitorSync(context.Background(), NewTaskPayload(TaskTypeCompetitorSync, "", "smartphone"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sync competitors failed")
	})
}

func TestHandler_HandleTrendAnalysis_WithLLM(t *testing.T) {
	logger := newTestLogger()

	t.Run("success", func(t *testing.T) {
		llm := &mockLLMClient{
			analyzeTrendsFunc: func(ctx context.Context, query string) (map[string]interface{}, error) {
				return map[string]interface{}{
					"trend": "upward",
					"opportunities": []string{"opp1", "opp2"},
				}, nil
			},
		}
		handler := NewHandler(logger, nil, llm)

		err := handler.HandleTrendAnalysis(context.Background(), NewTaskPayload(TaskTypeTrendAnalysis, "", "market trends"))
		assert.NoError(t, err)
	})

	t.Run("llm error", func(t *testing.T) {
		llm := &mockLLMClient{
			analyzeTrendsFunc: func(ctx context.Context, query string) (map[string]interface{}, error) {
				return nil, fmt.Errorf("llm error")
			},
		}
		handler := NewHandler(logger, nil, llm)

		err := handler.HandleTrendAnalysis(context.Background(), NewTaskPayload(TaskTypeTrendAnalysis, "", "market trends"))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "analyze trends failed")
	})
}

func TestHandler_HandleSupplyAlert_WithLLM(t *testing.T) {
	logger := newTestLogger()

	t.Run("success", func(t *testing.T) {
		llm := &mockLLMClient{
			checkSupplyChainFunc: func(ctx context.Context, productID string) ([]string, error) {
				return []string{"Risk1", "Risk2"}, nil
			},
		}
		handler := NewHandler(logger, nil, llm)

		err := handler.HandleSupplyAlert(context.Background(), NewTaskPayload(TaskTypeSupplyAlert, "prod-001", ""))
		assert.NoError(t, err)
	})

	t.Run("llm error", func(t *testing.T) {
		llm := &mockLLMClient{
			checkSupplyChainFunc: func(ctx context.Context, productID string) ([]string, error) {
				return nil, fmt.Errorf("llm error")
			},
		}
		handler := NewHandler(logger, nil, llm)

		err := handler.HandleSupplyAlert(context.Background(), NewTaskPayload(TaskTypeSupplyAlert, "prod-001", ""))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "check supply chain failed")
	})
}

func TestHandler_DevFallbacks(t *testing.T) {
	logger := newTestLogger()
	handler := NewHandler(logger, nil, nil) // no scraper, no llm

	t.Run("devScrapePrices returns fallback prices", func(t *testing.T) {
		prices, err := handler.devScrapePrices(context.Background(), "prod-001")
		assert.NoError(t, err)
		assert.NotEmpty(t, prices)
	})

	t.Run("devSyncCompetitorData returns fallback competitors", func(t *testing.T) {
		competitors, err := handler.devSyncCompetitorData(context.Background(), "query")
		assert.NoError(t, err)
		assert.NotEmpty(t, competitors)
	})
}
