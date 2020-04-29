package env

import (
	"os"
	"testing"

	"github.com/edaniszewski/chart-releaser/pkg/env"
	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/stretchr/testify/assert"
)

func TestStage_Name(t *testing.T) {
	assert.Equal(t, "env", Stage{}.Name())
}

func TestStage_String(t *testing.T) {
	assert.Equal(t, "loading environment variables", Stage{}.String())
}

func TestStage_Run_RepoGithub(t *testing.T) {
	if err := os.Setenv(env.GithubToken, "abc123"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Unsetenv(env.GithubToken); err != nil {
			t.Fatal(err)
		}
	}()

	context := ctx.Context{
		Repository: ctx.Repository{
			Type: ctx.RepoGithub,
		},
	}
	assert.Equal(t, "", context.Token)

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", context.Token)
}

func TestStage_Run_NoRepoType(t *testing.T) {
	context := ctx.Context{
		Repository: ctx.Repository{
			Type: "",
		},
	}
	assert.Equal(t, "", context.Token)

	err := Stage{}.Run(&context)
	assert.Error(t, err)
	assert.Equal(t, ErrRepoTypeNotSet, err)
	assert.Equal(t, "", context.Token)
}

func TestStage_Run_TokenNotSet(t *testing.T) {
	context := ctx.Context{
		Repository: ctx.Repository{
			Type: ctx.RepoGithub,
		},
	}
	assert.Equal(t, "", context.Token)

	err := Stage{}.Run(&context)
	assert.Error(t, err)
	assert.Equal(t, ErrGithubTokenNotSet, err)
	assert.Equal(t, "", context.Token)
}

func TestStage_Run_UnsupportedRepoType(t *testing.T) {
	context := ctx.Context{
		Repository: ctx.Repository{
			Type: ctx.RepoType("test-repo"),
		},
	}
	assert.Equal(t, "", context.Token)

	err := Stage{}.Run(&context)
	assert.Error(t, err)
	assert.Equal(t, ErrUnsupportedRepoType, err)
	assert.Equal(t, "", context.Token)
}

func TestStage_Run_TokenNotSetDryRun(t *testing.T) {
	context := ctx.Context{
		Repository: ctx.Repository{
			Type: ctx.RepoGithub,
		},
		DryRun: true,
	}
	assert.Equal(t, "", context.Token)

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
	assert.Equal(t, "", context.Token)
	assert.EqualError(t, context.Errors(), "\nErrors:\n • GITHUB_TOKEN environment variable not set\n\n")
}
