package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinExists(t *testing.T) {
	assert.True(t, BinExists("git"))
}

func TestBinNotExists(t *testing.T) {
	assert.False(t, BinExists("azksiejfng3las"))
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{
			"",
			"",
		},
		{
			"the quick brown fox jumps over the lazy dog",
			"the quick brown fox jumps over the lazy dog",
		},
		{
			"the quick\nbrown fox\njumps over\nthe lazy dog",
			"the quick brown fox jumps over the lazy dog",
		},
		{
			"the quick brown fox\r\njumps over the lazy dog",
			"the quick brown fox jumps over the lazy dog",
		},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			assert.Equal(t, test.expected, Normalize(test.in))
		})
	}
}
