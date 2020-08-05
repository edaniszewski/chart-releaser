package render

import (
	"testing"

	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/stretchr/testify/assert"
)

func TestStage_Name(t *testing.T) {
	assert.Equal(t, "render", Stage{}.Name())
}

func TestStage_String(t *testing.T) {
	assert.Equal(t, "rendering templates", Stage{}.String())
}

func TestStage_RunGitRefError(t *testing.T) {
	context := ctx.Context{
		Git: ctx.Git{
			Ref: "{{end}}",
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "template: git-ref:1: unexpected {{end}}")

	assert.Equal(t, "", context.Git.Ref)
	assert.Equal(t, "", context.Git.Base)
	assert.Equal(t, "", context.Release.ChartCommitMsg)
	assert.Equal(t, "", context.Release.PRTitle)
	assert.Equal(t, "", context.Release.PRBody)
}

func TestStage_RunGitBaseError(t *testing.T) {
	context := ctx.Context{
		Git: ctx.Git{
			Ref:  "ref-branch",
			Base: "{{end}}",
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "template: git-base:1: unexpected {{end}}")

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "", context.Git.Base)
	assert.Equal(t, "", context.Release.ChartCommitMsg)
	assert.Equal(t, "", context.Release.PRTitle)
	assert.Equal(t, "", context.Release.PRBody)
}

func TestStage_RunReleaseUpdateMsgError(t *testing.T) {
	context := ctx.Context{
		Git: ctx.Git{
			Ref:  "ref-branch",
			Base: "base-branch",
		},
		Release: ctx.Release{
			ChartCommitMsg: "{{end}}",
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "template: update-commit:1: unexpected {{end}}")

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "base-branch", context.Git.Base)
	assert.Equal(t, "", context.Release.ChartCommitMsg)
	assert.Equal(t, "", context.Release.PRTitle)
	assert.Equal(t, "", context.Release.PRBody)
}

func TestStage_RunReleasePRTitleError(t *testing.T) {
	context := ctx.Context{
		Git: ctx.Git{
			Ref:  "ref-branch",
			Base: "base-branch",
		},
		Release: ctx.Release{
			ChartCommitMsg: "update-msg",
			PRTitle:        "{{end}}",
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "template: pr-title:1: unexpected {{end}}")

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "base-branch", context.Git.Base)
	assert.Equal(t, "update-msg", context.Release.ChartCommitMsg)
	assert.Equal(t, "", context.Release.PRTitle)
	assert.Equal(t, "", context.Release.PRBody)
}

func TestStage_RunReleasePRBodyError(t *testing.T) {
	context := ctx.Context{
		Git: ctx.Git{
			Ref:  "ref-branch",
			Base: "base-branch",
		},
		Release: ctx.Release{
			ChartCommitMsg: "update-msg",
			PRTitle:        "pr-title",
			PRBody:         "{{end}}",
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "template: pr-body:1: unexpected {{end}}")

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "base-branch", context.Git.Base)
	assert.Equal(t, "update-msg", context.Release.ChartCommitMsg)
	assert.Equal(t, "pr-title", context.Release.PRTitle)
	assert.Equal(t, "", context.Release.PRBody)
}
