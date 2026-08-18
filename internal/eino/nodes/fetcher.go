package nodes

import (
	"context"
	"errors"
)

// FetcherInput represents input to the fetcher node.
type FetcherInput struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

// FetcherOutput represents output from the fetcher node.
type FetcherOutput struct {
	Data  string `json:"data"`
	Found bool   `json:"found"`
}

// ErrFetcherFailed is returned when fetcher fails to fetch data.
var ErrFetcherFailed = errors.New("fetcher failed to fetch data")

// FetcherNode fetches data from various sources.
type FetcherNode struct {
	// InjectError is used for testing error scenarios
	InjectError error
}

// NewFetcherNode creates a new FetcherNode.
func NewFetcherNode() *FetcherNode {
	return &FetcherNode{}
}

// Invoke executes the fetcher node.
func (n *FetcherNode) Invoke(ctx context.Context, input *FetcherInput) (*FetcherOutput, error) {
	if n.InjectError != nil {
		return nil, n.InjectError
	}
	return &FetcherOutput{
		Data:  "Fetched data for: " + input.Query,
		Found: true,
	}, nil
}
