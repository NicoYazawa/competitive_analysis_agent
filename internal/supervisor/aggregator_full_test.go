package supervisor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregator_generateSummary(t *testing.T) {
	agg := NewAggregator()

	tests := []struct {
		name     string
		result   *AggregatedResult
		contains string
	}{
		{
			name: "all results present",
			result: &AggregatedResult{
				MarketTrend:       &MarketTrendResult{Trend: "Trend A"},
				CompetitorAnalysis: &CompetitorResult{Competitors: []CompetitorInsight{{Name: "C1"}}},
				PricingStrategy:   &PricingResult{RecommendedPriceRange: "$10-$20"},
				SupplyChainInfo:   &SupplyChainResult{RiskLevel: "Low"},
			},
			contains: "Trend A",
		},
		{
			name:     "empty result",
			result:   &AggregatedResult{},
			contains: "No analysis results available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := agg.generateSummary(tt.result)
			assert.Contains(t, summary, tt.contains)
		})
	}
}

func TestAggregator_Aggregate_WithErrors(t *testing.T) {
	agg := NewAggregator()

	tasks := []*Task{
		{
			Type:   TaskTypeMarketTrend,
			Result: nil,
		},
		{
			Type:   TaskTypeCompetitor,
			Result: "error: something went wrong",
		},
	}

	result := agg.Aggregate(tasks)

	assert.Equal(t, 2, result.TaskCount)
	assert.Nil(t, result.MarketTrend)
}

func TestAggregator_Aggregate_SupplyChainOnly(t *testing.T) {
	agg := NewAggregator()

	tasks := []*Task{
		{
			Type: TaskTypeSupplyChain,
			Result: &SupplyChainResult{
				Status:    "Stable",
				RiskLevel: "Low",
				Factors:   []string{"Factor 1", "Factor 2"},
			},
		},
	}

	result := agg.Aggregate(tasks)

	assert.Equal(t, 1, result.TaskCount)
	assert.NotNil(t, result.SupplyChainInfo)
	assert.Equal(t, "Low", result.SupplyChainInfo.RiskLevel)
	assert.Contains(t, result.Summary, "Supply Chain")
}

func TestJoinStrings(t *testing.T) {
	assert.Equal(t, "", joinStrings([]string{}, ","))
	assert.Equal(t, "a", joinStrings([]string{"a"}, ","))
	assert.Equal(t, "a,b,c", joinStrings([]string{"a", "b", "c"}, ","))
}
