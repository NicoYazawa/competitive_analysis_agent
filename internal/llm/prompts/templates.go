package prompts

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateData represents data for prompt templates.
type TemplateData struct {
	Query          string
	Competitors    []CompetitorInfo
	PriceHistory   []PricePoint
	MarketTrend    string
	SupplyChainInfo string
}

// CompetitorInfo holds competitor data for templates.
type CompetitorInfo struct {
	Name         string
	Platform     string
	Price        float64
	Rating       float64
	ReviewCount  int
}

// PricePoint represents a price data point.
type PricePoint struct {
	Timestamp string
	Price     float64
	Platform  string
}

// RenderTrendAnalysis renders the trend analysis prompt.
func RenderTrendAnalysis(data *TemplateData) (string, error) {
	tmpl := `You are a market trend analyst. Analyze the following market data and provide insights.

Query: {{.Query}}

Market Trend Summary:
{{.MarketTrend}}

Recent Competitor Data:
{{range .Competitors}}
- {{.Name}} on {{.Platform}}: ${{.Price}} (Rating: {{.Rating}}, Reviews: {{.ReviewCount}})
{{end}}

Provide a concise trend analysis including:
1. Price positioning opportunities
2. Market demand signals
3. Competitive landscape summary
`
	return renderTemplate(tmpl, data)
}

// RenderCompetitorAnalysis renders the competitor analysis prompt.
func RenderCompetitorAnalysis(data *TemplateData) (string, error) {
	tmpl := `You are a competitor analysis expert. Analyze the following competitor data.

Query: {{.Query}}

Competitors:
{{range .Competitors}}
- {{.Name}} ({{.Platform}})
  Price: ${{.Price}}
  Rating: {{.Rating}}/5 ({{.ReviewCount}} reviews)
{{end}}

Provide:
1. Strengths and weaknesses of each competitor
2. Pricing strategy analysis
3. Market positioning recommendations
`
	return renderTemplate(tmpl, data)
}

// RenderPricingStrategy renders the pricing strategy prompt.
func RenderPricingStrategy(data *TemplateData) (string, error) {
	tmpl := `You are a pricing strategy expert. Based on competitor and market data, recommend optimal pricing.

Query: {{.Query}}

Price History:
{{range .PriceHistory}}
- {{.Timestamp}}: ${{.Price}} on {{.Platform}}
{{end}}

Competitor Prices:
{{range .Competitors}}
- {{.Name}}: ${{.Price}}
{{end}}

{{if .SupplyChainInfo}}
Supply Chain Info:
{{.SupplyChainInfo}}
{{end}}

Provide:
1. Recommended price range
2. Pricing strategy rationale
3. Competitive positioning advice
`
	return renderTemplate(tmpl, data)
}

func renderTemplate(tmplStr string, data *TemplateData) (string, error) {
	tmpl, err := template.New("prompt").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
