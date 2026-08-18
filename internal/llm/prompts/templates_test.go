package prompts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderTrendAnalysis(t *testing.T) {
	data := &TemplateData{
		Query:       "Analyze smartphone market trends",
		MarketTrend: "Growing demand for mid-range devices",
		Competitors: []CompetitorInfo{
			{Name: "Samsung", Platform: "Amazon", Price: 599.99, Rating: 4.5, ReviewCount: 1000},
		},
	}

	result, err := RenderTrendAnalysis(data)

	require.NoError(t, err)
	assert.Contains(t, result, "smartphone market trends")
	assert.Contains(t, result, "Growing demand for mid-range devices")
	assert.Contains(t, result, "Samsung")
}

func TestRenderCompetitorAnalysis(t *testing.T) {
	data := &TemplateData{
		Query: "Compare laptop competitors",
		Competitors: []CompetitorInfo{
			{Name: "Dell", Platform: "Amazon", Price: 999.99, Rating: 4.2, ReviewCount: 500},
		},
	}

	result, err := RenderCompetitorAnalysis(data)

	require.NoError(t, err)
	assert.Contains(t, result, "laptop competitors")
	assert.Contains(t, result, "Dell")
}

func TestRenderPricingStrategy(t *testing.T) {
	data := &TemplateData{
		Query: "Recommend pricing for new product",
		PriceHistory: []PricePoint{
			{Timestamp: "2024-01-01", Price: 99.99, Platform: "Amazon"},
		},
		Competitors: []CompetitorInfo{
			{Name: "Competitor A", Price: 89.99},
		},
		SupplyChainInfo: "Stable supply chain",
	}

	result, err := RenderPricingStrategy(data)

	require.NoError(t, err)
	assert.Contains(t, result, "new product")
	assert.Contains(t, result, "99.99")
	assert.Contains(t, result, "Stable supply chain")
}

func TestRenderPricingStrategy_WithoutSupplyChain(t *testing.T) {
	data := &TemplateData{
		Query: "Recommend pricing for new product",
		PriceHistory: []PricePoint{
			{Timestamp: "2024-01-01", Price: 99.99, Platform: "Amazon"},
		},
		Competitors: []CompetitorInfo{
			{Name: "Competitor A", Price: 89.99},
		},
	}

	result, err := RenderPricingStrategy(data)

	require.NoError(t, err)
	assert.NotContains(t, result, "Supply Chain Info:")
}
