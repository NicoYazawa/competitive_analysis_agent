package valueobject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScore(t *testing.T) {
	score := NewScore(85.5, 100, "Performance", "rating")

	assert.Equal(t, 85.5, score.Value)
	assert.Equal(t, 100.0, score.MaxValue)
	assert.Equal(t, "Performance", score.Label)
	assert.Equal(t, "rating", score.Category)
}

func TestScorePercentage(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		maxValue float64
		expected float64
	}{
		{"full", 100, 100, 100},
		{"half", 50, 100, 50},
		{"quarter", 25, 100, 25},
		{"zero", 0, 100, 0},
		{"zero max", 50, 0, 0},
		{"decimal", 33.33, 100, 33.33},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := NewScore(tt.value, tt.maxValue, "Test", "test")
			assert.Equal(t, tt.expected, score.Percentage())
		})
	}
}

func TestScorePercentageEdgeCases(t *testing.T) {
	// Test zero max value (should not panic)
	score := NewScore(50, 0, "Test", "test")
	assert.Equal(t, 0.0, score.Percentage())

	// Test very small values
	score = NewScore(0.001, 100, "Test", "test")
	assert.InDelta(t, 0.001, score.Percentage(), 0.0001)
}
