package render

import (
	context "github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/edaniszewski/chart-releaser/pkg/v1/utils"
)

// Stage for the "render" step of the update pipeline.
type Stage struct{}

// Name of the stage.
func (Stage) Name() string {
	return "render"
}

// String describes what the stage does.
func (Stage) String() string {
	return "rendering templates"
}

// Run the operations defined for the stage.
func (Stage) Run(ctx *context.Context) error {
	var err error

	ctx.Git.Ref, err = utils.RenderTemplate(ctx, "git-ref", ctx.Git.Ref)
	if err != nil {
		return err
	}

	ctx.Git.Base, err = utils.RenderTemplate(ctx, "git-base", ctx.Git.Base)
	if err != nil {
		return err
	}

	ctx.Release.UpdateCommitMsg, err = utils.RenderTemplate(ctx, "update-commit", ctx.Release.UpdateCommitMsg)
	if err != nil {
		return err
	}

	ctx.Release.PRTitle, err = utils.RenderTemplate(ctx, "pr-title", ctx.Release.PRTitle)
	if err != nil {
		return err
	}

	ctx.Release.PRBody, err = utils.RenderTemplate(ctx, "pr-body", ctx.Release.PRBody)
	if err != nil {
		return err
	}

	return nil
}
