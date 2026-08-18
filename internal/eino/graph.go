package eino

import (
	"context"

	"competitive-analysis-agent/internal/eino/nodes"
	"competitive-analysis-agent/internal/llm"
)

// Graph represents the Eino DAG for agent orchestration.
type Graph struct {
	fetcher    *nodes.FetcherNode
	analyzer   *nodes.AnalyzerNode
	synthesizer *nodes.SynthesizerNode
}

// NewGraph creates a new Eino Graph.
func NewGraph(llmClient llm.ChatModel) *Graph {
	return &Graph{
		fetcher:    nodes.NewFetcherNode(),
		analyzer:   nodes.NewAnalyzerNode(llmClient),
		synthesizer: nodes.NewSynthesizerNode(),
	}
}

// ExecuteResult holds the result of graph execution.
type ExecuteResult struct {
	Strategy        string   `json:"strategy"`
	Recommendations []string `json:"recommendations"`
	AnalysisResults []string `json:"analysis_results"`
}

// Execute runs the DAG: fetcher -> analyzer -> synthesizer.
func (g *Graph) Execute(ctx context.Context, query string) (*ExecuteResult, error) {
	// Step 1: Fetch data
	fetcherOut, err := g.fetcher.Invoke(ctx, &nodes.FetcherInput{
		Query: query,
		Type:  "market_data",
	})
	if err != nil {
		return nil, err
	}

	// Step 2: Analyze data
	analyzerOut, err := g.analyzer.Invoke(ctx, &nodes.AnalyzerInput{
		Query:    query,
		RawData:  fetcherOut.Data,
	})
	if err != nil {
		return nil, err
	}

	// Step 3: Synthesize results
	synthOut, err := g.synthesizer.Invoke(ctx, &nodes.SynthesizerInput{
		AnalysisResults: []string{analyzerOut.Analysis},
		Query:           query,
	})
	if err != nil {
		return nil, err
	}

	return &ExecuteResult{
		Strategy:        synthOut.FinalStrategy,
		Recommendations: synthOut.Recommendations,
		AnalysisResults: []string{analyzerOut.Analysis},
	}, nil
}
