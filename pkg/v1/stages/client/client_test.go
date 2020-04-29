package client

import (
	"testing"

	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/stretchr/testify/assert"
)

func TestStage_Name(t *testing.T) {
	assert.Equal(t, "client", Stage{}.Name())
}

func TestStage_String(t *testing.T) {
	assert.Equal(t, "creating repository client", Stage{}.String())
}

func TestStage_Run(t *testing.T) {
	context := ctx.Context{
		Token: "test-token",
		Repository: ctx.Repository{
			Type: ctx.RepoGithub,
		},
	}
	assert.Nil(t, context.Client)

	err := Stage{}.Run(&context)
	assert.Nil(t, err)
	assert.NotNil(t, context.Client)
}

func TestStage_Run_ErrorNoToken(t *testing.T) {
	context := ctx.Context{
		Repository: ctx.Repository{
			Type: ctx.RepoGithub,
		},
	}
	assert.Nil(t, context.Client)

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "no token provided to github client")
	assert.Nil(t, context.Client)
}

func TestStage_Run_ErrorRepoTypeNotSet(t *testing.T) {
	context := ctx.Context{
		Repository: ctx.Repository{},
	}
	assert.Nil(t, context.Client)

	err := Stage{}.Run(&context)
	assert.Equal(t, ErrRepoTypeNotSet, err)
	assert.Nil(t, context.Client)
}

func TestStage_Run_ErrorUnsupportedRepoType(t *testing.T) {
	context := ctx.Context{
		Repository: ctx.Repository{
			Type: ctx.RepoType("unsupported-repo"),
		},
	}
	assert.Nil(t, context.Client)

	err := Stage{}.Run(&context)
	assert.Equal(t, ErrUnsupportedRepoType, err)
	assert.Nil(t, context.Client)
}
