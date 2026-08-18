package supervisor

import (
	"fmt"
)

// AggregatedResult represents the final aggregated result from all agents.
type AggregatedResult struct {
	Summary           string                 `json:"summary"`
	MarketTrend       *MarketTrendResult     `json:"market_trend,omitempty"`
	CompetitorAnalysis *CompetitorResult     `json:"competitor_analysis,omitempty"`
	PricingStrategy   *PricingResult         `json:"pricing_strategy,omitempty"`
	SupplyChainInfo   *SupplyChainResult     `json:"supply_chain,omitempty"`
	TaskCount         int                    `json:"task_count"`
	Errors            []string               `json:"errors,omitempty"`
}

// MarketTrendResult holds market trend analysis result.
type MarketTrendResult struct {
	Trend        string   `json:"trend"`
	Opportunities []string `json:"opportunities"`
	DemandSignal string   `json:"demand_signal"`
}

// CompetitorResult holds competitor analysis result.
type CompetitorResult struct {
	Analysis    string                `json:"analysis"`
	Competitors []CompetitorInsight   `json:"competitors"`
}

// CompetitorInsight holds insight about a single competitor.
type CompetitorInsight struct {
	Name       string `json:"name"`
	Strength   string `json:"strength"`
	Weakness   string `json:"weakness"`
	Strategy   string `json:"strategy"`
}

// PricingResult holds pricing strategy result.
type PricingResult struct {
	RecommendedPriceRange string   `json:"recommended_price_range"`
	Rationale             string   `json:"rationale"`
	Positioning           string   `json:"positioning"`
}

// SupplyChainResult holds supply chain analysis result.
type SupplyChainResult struct {
	Status     string   `json:"status"`
	RiskLevel  string   `json:"risk_level"`
	Factors    []string `json:"factors"`
}

// Aggregator combines results from multiple agents into a final result.
type Aggregator struct{}

// NewAggregator creates a new Aggregator.
func NewAggregator() *Aggregator {
	return &Aggregator{}
}

// Aggregate combines task results into a final result.
func (a *Aggregator) Aggregate(tasks []*Task) *AggregatedResult {
	result := &AggregatedResult{
		TaskCount: len(tasks),
		Errors:    []string{},
	}

	for _, task := range tasks {
		if task.Result == nil {
			continue
		}

		switch task.Type {
		case TaskTypeMarketTrend:
			if mr, ok := task.Result.(*MarketTrendResult); ok {
				result.MarketTrend = mr
			}
		case TaskTypeCompetitor:
			if cr, ok := task.Result.(*CompetitorResult); ok {
				result.CompetitorAnalysis = cr
			}
		case TaskTypePricing:
			if pr, ok := task.Result.(*PricingResult); ok {
				result.PricingStrategy = pr
			}
		case TaskTypeSupplyChain:
			if sr, ok := task.Result.(*SupplyChainResult); ok {
				result.SupplyChainInfo = sr
			}
		}

		// Check for errors
		if errStr, isErr := task.Result.(string); isErr && len(errStr) > 0 && errStr[0] == 'e' {
			// Simple check for error string prefix
			if len(errStr) > 5 && errStr[:5] == "error" {
				result.Errors = append(result.Errors, errStr)
			}
		}
	}

	// Generate summary
	result.Summary = a.generateSummary(result)

	return result
}

func (a *Aggregator) generateSummary(r *AggregatedResult) string {
	parts := []string{}

	if r.MarketTrend != nil {
		parts = append(parts, fmt.Sprintf("Market Trend: %s", r.MarketTrend.Trend))
	}
	if r.CompetitorAnalysis != nil {
		parts = append(parts, fmt.Sprintf("Competitor Analysis completed for %d competitors", len(r.CompetitorAnalysis.Competitors)))
	}
	if r.PricingStrategy != nil {
		parts = append(parts, fmt.Sprintf("Pricing: %s", r.PricingStrategy.RecommendedPriceRange))
	}
	if r.SupplyChainInfo != nil {
		parts = append(parts, fmt.Sprintf("Supply Chain: %s risk", r.SupplyChainInfo.RiskLevel))
	}

	if len(parts) == 0 {
		return "No analysis results available"
	}

	return fmt.Sprintf("Analysis complete. %s.", joinStrings(parts, " "))
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
