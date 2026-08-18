package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float64
	}{
		{
			name:     "identical vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "opposite vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{-1, 0, 0},
			expected: -1.0,
		},
		{
			name:     "perpendicular vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "normalized vectors",
			a:        []float32{0.707, 0.707, 0},
			b:        []float32{0.707, 0.707, 0},
			expected: 1.0,
		},
		{
			name:     "different length vectors",
			a:        []float32{1, 0},
			b:        []float32{1, 0, 0},
			expected: 0.0,
		},
		{
			name:     "zero vector a",
			a:        []float32{0, 0, 0},
			b:        []float32{1, 0, 0},
			expected: 0.0,
		},
		{
			name:     "zero vector b",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 0, 0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CosineSimilarity(tt.a, tt.b)
			assert.InDelta(t, tt.expected, result, 0.001)
		})
	}
}

func TestBytesToFloat32(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []float32
	}{
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: nil,
		},
		{
			name:     "nil bytes",
			input:    nil,
			expected: nil,
		},
		{
			name:     "valid json array",
			input:    []byte(`[1.0, 2.0, 3.0]`),
			expected: []float32{1.0, 2.0, 3.0},
		},
		{
			name:     "invalid json",
			input:    []byte(`not json`),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytesToFloat32(tt.input)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestJSONMap(t *testing.T) {
	m := JSONMap{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	assert.Equal(t, "value1", m["key1"])
	assert.Equal(t, 123, m["key2"])
	assert.Equal(t, true, m["key3"])
}

func TestVectorEmbedding(t *testing.T) {
	emb := VectorEmbedding{
		ID:       "test-id",
		Entity:   "competitor",
		Data:     `{"name": "Test"}`,
		Vector:   []float32{0.1, 0.2, 0.3},
		Metadata: JSONMap{"source": "test"},
	}

	assert.Equal(t, "test-id", emb.ID)
	assert.Equal(t, "competitor", emb.Entity)
	assert.Equal(t, []float32{0.1, 0.2, 0.3}, emb.Vector)
	assert.Equal(t, "test", emb.Metadata["source"])
}
