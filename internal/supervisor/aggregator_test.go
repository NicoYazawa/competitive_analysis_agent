package supervisor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregator_Aggregate(t *testing.T) {
	agg := NewAggregator()

	tasks := []*Task{
		{
			Type: TaskTypeMarketTrend,
			Result: &MarketTrendResult{
				Trend:         "Bullish",
				Opportunities: []string{"Opp1"},
				DemandSignal:  "Strong",
			},
		},
		{
			Type: TaskTypeCompetitor,
			Result: &CompetitorResult{
				Analysis: "Competitive analysis content",
				Competitors: []CompetitorInsight{
					{Name: "Comp A", Strength: "Brand", Weakness: "Price", Strategy: "Premium"},
				},
			},
		},
		{
			Type: TaskTypePricing,
			Result: &PricingResult{
				RecommendedPriceRange: "$50-$60",
				Rationale:             "Based on competition",
				Positioning:           "Mid-range",
			},
		},
	}

	result := agg.Aggregate(tasks)

	assert.Equal(t, 3, result.TaskCount)
	assert.NotNil(t, result.MarketTrend)
	assert.NotNil(t, result.CompetitorAnalysis)
	assert.NotNil(t, result.PricingStrategy)
	assert.Contains(t, result.Summary, "Bullish")
	assert.Contains(t, result.Summary, "$50-$60")
}

func TestAggregator_Aggregate_EmptyTasks(t *testing.T) {
	agg := NewAggregator()

	result := agg.Aggregate([]*Task{})

	assert.Equal(t, 0, result.TaskCount)
	assert.Equal(t, "No analysis results available", result.Summary)
}

func TestAggregator_Aggregate_MissingResults(t *testing.T) {
	agg := NewAggregator()

	tasks := []*Task{
		{
			Type:   TaskTypeMarketTrend,
			Result: nil,
		},
	}

	result := agg.Aggregate(tasks)

	assert.Equal(t, 1, result.TaskCount)
	assert.Nil(t, result.MarketTrend)
}
