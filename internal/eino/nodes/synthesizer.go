package nodes

import (
	"context"
	"errors"
	"fmt"
)

// SynthesizerInput represents input to the synthesizer node.
type SynthesizerInput struct {
	AnalysisResults []string `json:"analysis_results"`
	Query           string   `json:"query"`
}

// SynthesizerOutput represents output from the synthesizer node.
type SynthesizerOutput struct {
	FinalStrategy   string   `json:"final_strategy"`
	Recommendations []string `json:"recommendations"`
}

// ErrSynthesisFailed is returned when synthesis fails.
var ErrSynthesisFailed = errors.New("synthesis failed")

// SynthesizerNode synthesizes multiple analysis results into a final strategy.
type SynthesizerNode struct {
	// InjectError is used for testing error scenarios
	InjectError error
}

// NewSynthesizerNode creates a new SynthesizerNode.
func NewSynthesizerNode() *SynthesizerNode {
	return &SynthesizerNode{}
}

// Invoke executes the synthesizer node.
func (n *SynthesizerNode) Invoke(ctx context.Context, input *SynthesizerInput) (*SynthesizerOutput, error) {
	if n.InjectError != nil {
		return nil, n.InjectError
	}
	if len(input.AnalysisResults) == 0 {
		return &SynthesizerOutput{
			FinalStrategy:   "No analysis results to synthesize",
			Recommendations: []string{},
		}, nil
	}

	// Simple synthesis logic
	synthesis := fmt.Sprintf("Based on %d analysis results for: %s", len(input.AnalysisResults), input.Query)

	return &SynthesizerOutput{
		FinalStrategy:   synthesis,
		Recommendations: []string{"Recommendation 1", "Recommendation 2", "Recommendation 3"},
	}, nil
}
