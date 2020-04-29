package strategies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPublishStrategies(t *testing.T) {
	strategies := ListPublishStrategies()
	assert.Len(t, strategies, 2)
	assert.Equal(t, PublishCommit, strategies[0])
	assert.Equal(t, PublishPullRequest, strategies[1])
}

func TestPublishStrategyFromString(t *testing.T) {
	tests := []struct {
		str      string
		expected PublishStrategy
	}{
		{
			str:      "commit",
			expected: PublishCommit,
		},
		{
			str:      "Commit",
			expected: PublishCommit,
		},
		{
			str:      "COMMIT",
			expected: PublishCommit,
		},
		{
			str:      "pull request",
			expected: PublishPullRequest,
		},
		{
			str:      "Pull Request",
			expected: PublishPullRequest,
		},
		{
			str:      "PULL REQUEST",
			expected: PublishPullRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			actual, err := PublishStrategyFromString(test.str)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestPublishStrategyFromString_Error(t *testing.T) {
	strategy, err := PublishStrategyFromString("test-string")
	assert.Error(t, err)
	assert.Empty(t, strategy)
}
