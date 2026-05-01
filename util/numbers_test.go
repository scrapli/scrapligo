package util_test

import (
	"math"
	"testing"

	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
)

func TestSafeInt64ToUint64(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected uint64
	}{
		{
			name:     "positive value",
			input:    100,
			expected: 100,
		},
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
		{
			name:     "negative value",
			input:    -1,
			expected: math.MaxUint64,
		},
		{
			name:     "large negative value",
			input:    math.MinInt64,
			expected: math.MaxUint64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := scrapligoutil.SafeInt64ToUint64(tt.input)
			if actual != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, actual)
			}
		})
	}
}

func TestSafeUint64ToInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected int64
	}{
		{
			name:     "positive value",
			input:    100,
			expected: 100,
		},
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
		{
			name:     "max int64",
			input:    uint64(math.MaxInt64),
			expected: math.MaxInt64,
		},
		{
			name:     "overflowing int64",
			input:    uint64(math.MaxInt64) + 1,
			expected: math.MaxInt64,
		},
		{
			name:     "max uint64",
			input:    math.MaxUint64,
			expected: math.MaxInt64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := scrapligoutil.SafeUint64ToInt64(tt.input)
			if actual != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, actual)
			}
		})
	}
}
