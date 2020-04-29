package utils

import (
	"testing"

	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/stretchr/testify/assert"
)

func TestParseRepository(t *testing.T) {
	tests := []struct {
		repo     string
		expected ctx.Repository
	}{
		{
			repo: "test/charts",
			expected: ctx.Repository{
				Type:  ctx.RepoGithub,
				Owner: "test",
				Name:  "charts",
			},
		},
		{
			repo: "github/test/charts",
			expected: ctx.Repository{
				Type:  ctx.RepoGithub,
				Owner: "test",
				Name:  "charts",
			},
		},
		{
			repo: "github.com/test/charts",
			expected: ctx.Repository{
				Type:  ctx.RepoGithub,
				Owner: "test",
				Name:  "charts",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.repo, func(t *testing.T) {
			actual, err := ParseRepository(test.repo)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestParseRepository_ErrorNumberParts(t *testing.T) {
	_, err := ParseRepository("one/two/three/four")
	assert.Error(t, err)
}

func TestParseRepository_ErrorNumberParts2(t *testing.T) {
	_, err := ParseRepository("one")
	assert.Error(t, err)
}

func TestParseRepository_ErrorRepoType(t *testing.T) {
	_, err := ParseRepository("unsupportedrepo.com/test/charts")
	assert.Error(t, err)
}
