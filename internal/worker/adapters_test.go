package worker

import (
	"context"
	"errors"
	"testing"

	"competitive-analysis-agent/internal/llm"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// LLMAnalyzerAdapter tests - uses llm.ChatModel interface so can be mocked
// ---------------------------------------------------------------------------

func TestLLMAnalyzerAdapter_AnalyzeTrends(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLLM := &llm.MockLLM{
			RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
				return &llm.ChatResponse{Content: "upward trend"}, nil
			},
		}
		adapter := NewLLMAnalyzerAdapter(mockLLM)

		result, err := adapter.AnalyzeTrends(context.Background(), "market trends")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "upward trend", result["trend"])
		assert.NotEmpty(t, result["opportunities"])
		assert.NotEmpty(t, result["demand_signal"])
	})

	t.Run("llm error", func(t *testing.T) {
		mockLLM := &llm.MockLLM{
			ErrFn: func(ctx context.Context, req *llm.ChatRequest) error {
				return errors.New("llm error")
			},
		}
		adapter := NewLLMAnalyzerAdapter(mockLLM)

		result, err := adapter.AnalyzeTrends(context.Background(), "market trends")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestLLMAnalyzerAdapter_CheckSupplyChain(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockLLM := &llm.MockLLM{
			RespFn: func(ctx context.Context, req *llm.ChatRequest) (*llm.ChatResponse, error) {
				return &llm.ChatResponse{Content: "shipping risk"}, nil
			},
		}
		adapter := NewLLMAnalyzerAdapter(mockLLM)

		result, err := adapter.CheckSupplyChain(context.Background(), "prod-123")
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		// Default response includes the LLM content plus two extra risks
		assert.Contains(t, result, "shipping risk")
	})

	t.Run("llm error", func(t *testing.T) {
		mockLLM := &llm.MockLLM{
			ErrFn: func(ctx context.Context, req *llm.ChatRequest) error {
				return errors.New("llm error")
			},
		}
		adapter := NewLLMAnalyzerAdapter(mockLLM)

		result, err := adapter.CheckSupplyChain(context.Background(), "prod-123")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestNewLLMAnalyzerAdapter(t *testing.T) {
	mockLLM := &llm.MockLLM{}
	adapter := NewLLMAnalyzerAdapter(mockLLM)
	assert.NotNil(t, adapter)
}

// ---------------------------------------------------------------------------
// CollyScraperAdapter tests - requires concrete *scraper.CollyScraper type
// ---------------------------------------------------------------------------

// Note: CollyScraperAdapter and MultiPlatformAdapter use concrete struct types
// (*scraper.CollyScraper, *platforms.MultiPlatformScraper) rather than interfaces,
// making them difficult to unit test without refactoring to use interfaces.
// They are tested via integration tests instead.

func TestCollyScraperAdapter_RequiresRealScraper(t *testing.T) {
	t.Skip("CollyScraperAdapter depends on concrete scraper.CollyScraper type - tested via integration")
}

func TestNewCollyScraperAdapter_RequiresRealScraper(t *testing.T) {
	t.Skip("CollyScraperAdapter depends on concrete scraper.CollyScraper type - tested via integration")
}

func TestMultiPlatformAdapter_RequiresRealScraper(t *testing.T) {
	t.Skip("MultiPlatformAdapter depends on concrete platforms.MultiPlatformScraper type - tested via integration")
}

func TestScrapeCompetitorPricesViaScraper_RequiresRealScraper(t *testing.T) {
	t.Skip("ScrapeCompetitorPricesViaScraper depends on concrete scraper.CollyScraper type - tested via integration")
}
