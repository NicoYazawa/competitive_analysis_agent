package prompts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderTrendAnalysis_EmptyCompetitors(t *testing.T) {
	data := &TemplateData{
		Query:       "test query",
		MarketTrend: "trending",
		Competitors: []CompetitorInfo{},
	}

	result, err := RenderTrendAnalysis(data)

	require.NoError(t, err)
	assert.Contains(t, result, "test query")
	assert.Contains(t, result, "trending")
}

func TestRenderTrendAnalysis_MultipleCompetitors(t *testing.T) {
	data := &TemplateData{
		Query:       "smartphone analysis",
		MarketTrend: "growing",
		Competitors: []CompetitorInfo{
			{Name: "Samsung", Platform: "Amazon", Price: 599.99, Rating: 4.5, ReviewCount: 1000},
			{Name: "Apple", Platform: "Amazon", Price: 999.99, Rating: 4.8, ReviewCount: 2000},
			{Name: "Xiaomi", Platform: "Alibaba", Price: 299.99, Rating: 4.2, ReviewCount: 500},
		},
	}

	result, err := RenderTrendAnalysis(data)

	require.NoError(t, err)
	assert.Contains(t, result, "Samsung")
	assert.Contains(t, result, "Apple")
	assert.Contains(t, result, "Xiaomi")
	assert.Contains(t, result, "599.99")
	assert.Contains(t, result, "999.99")
}

func TestRenderCompetitorAnalysis_MultipleCompetitors(t *testing.T) {
	data := &TemplateData{
		Query: "laptop comparison",
		Competitors: []CompetitorInfo{
			{Name: "Dell", Platform: "Amazon", Price: 999.99, Rating: 4.2, ReviewCount: 500},
			{Name: "HP", Platform: "Amazon", Price: 899.99, Rating: 4.0, ReviewCount: 400},
			{Name: "Lenovo", Platform: "BestBuy", Price: 1099.99, Rating: 4.6, ReviewCount: 300},
		},
	}

	result, err := RenderCompetitorAnalysis(data)

	require.NoError(t, err)
	assert.Contains(t, result, "Dell")
	assert.Contains(t, result, "HP")
	assert.Contains(t, result, "Lenovo")
	assert.Contains(t, result, "999.99")
}

func TestRenderPricingStrategy_MultiplePriceHistory(t *testing.T) {
	data := &TemplateData{
		Query: "pricing recommendation",
		PriceHistory: []PricePoint{
			{Timestamp: "2024-01-01", Price: 99.99, Platform: "Amazon"},
			{Timestamp: "2024-02-01", Price: 95.99, Platform: "Amazon"},
			{Timestamp: "2024-03-01", Price: 89.99, Platform: "Amazon"},
		},
		Competitors: []CompetitorInfo{
			{Name: "CompA", Price: 85.00},
			{Name: "CompB", Price: 95.00},
		},
	}

	result, err := RenderPricingStrategy(data)

	require.NoError(t, err)
	assert.Contains(t, result, "99.99")
	assert.Contains(t, result, "95.99")
	assert.Contains(t, result, "89.99")
	assert.Contains(t, result, "CompA")
	assert.Contains(t, result, "CompB")
}

func TestRenderPricingStrategy_WithSupplyChain(t *testing.T) {
	data := &TemplateData{
		Query:         "pricing with supply chain",
		SupplyChainInfo: "Shipping costs increasing, supplier lead time 4 weeks",
	}

	result, err := RenderPricingStrategy(data)

	require.NoError(t, err)
	assert.Contains(t, result, "Shipping costs increasing")
}

func TestRenderTemplate_Invalid(t *testing.T) {
	// This tests the text/template parsing
	data := &TemplateData{Query: "test"}

	// Use an invalid template to test error handling
	_, err := renderTemplate("{{.InvalidField}}", data)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute template")
}

func TestTemplateData_Structure(t *testing.T) {
	data := TemplateData{
		Query:           "test",
		MarketTrend:     "trend",
		SupplyChainInfo: "supply",
		Competitors: []CompetitorInfo{
			{Name: "C1", Platform: "P1", Price: 10.0, Rating: 4.0, ReviewCount: 100},
		},
		PriceHistory: []PricePoint{
			{Timestamp: "2024-01-01", Price: 10.0, Platform: "Amazon"},
		},
	}

	assert.Equal(t, "test", data.Query)
	assert.Equal(t, "trend", data.MarketTrend)
	assert.Equal(t, "supply", data.SupplyChainInfo)
	assert.Len(t, data.Competitors, 1)
	assert.Len(t, data.PriceHistory, 1)
}

func TestCompetitorInfo_Structure(t *testing.T) {
	info := CompetitorInfo{
		Name:        "TestComp",
		Platform:    "Amazon",
		Price:       99.99,
		Rating:      4.5,
		ReviewCount: 1000,
	}

	assert.Equal(t, "TestComp", info.Name)
	assert.Equal(t, "Amazon", info.Platform)
	assert.Equal(t, 99.99, info.Price)
	assert.Equal(t, 4.5, info.Rating)
	assert.Equal(t, 1000, info.ReviewCount)
}

func TestPricePoint_Structure(t *testing.T) {
	pp := PricePoint{
		Timestamp: "2024-01-01",
		Price:     50.00,
		Platform:  "Alibaba",
	}

	assert.Equal(t, "2024-01-01", pp.Timestamp)
	assert.Equal(t, 50.00, pp.Price)
	assert.Equal(t, "Alibaba", pp.Platform)
}
