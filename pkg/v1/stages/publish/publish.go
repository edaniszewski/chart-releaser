package publish

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/edaniszewski/chart-releaser/pkg/client"

	"github.com/apex/log"
	"github.com/davecgh/go-spew/spew"
	"github.com/edaniszewski/chart-releaser/pkg/strategies"
	"github.com/edaniszewski/chart-releaser/pkg/templates"
	context "github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/edaniszewski/chart-releaser/pkg/v1/utils"
)

// Errors for the publish stage.
var (
	ErrUnsupportedPublishStrategy = errors.New("unsupported publish strategy specified")
	ErrNoChartChanges             = errors.New("chart file has no changes")
)

// Stage for the "publish" step of the update pipeline.
type Stage struct{}

// Name of the stage.
func (Stage) Name() string {
	return "publish"
}

// String describes what the stage does.
func (Stage) String() string {
	return "publishing changes"
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

	log.Debugf("chart-release context:\n%v", spew.Sdump(ctx))

	// Check that the tag matches the release constraints.
	for _, m := range ctx.Release.Matches {
		if !m.MatchString(ctx.Git.Tag) {
			if ctx.DryRun {
				log.Warnf("dry-run: tag release (%s) does not match release constraint '%s'", ctx.Git.Tag, m.String())
				continue
			} else {
				log.Infof("tag release (%s) does not match release constraint '%s': will not update", ctx.Git.Tag, m.String())
				return nil
			}
		}
	}

	// Check that the tag does not match an the ignore constraints.
	for _, i := range ctx.Release.Ignores {
		if i.MatchString(ctx.Git.Tag) {
			if ctx.DryRun {
				log.Warnf("dry-run: tag release (%s) matches release ignore constraint '%s'", ctx.Git.Tag, i.String())
				continue
			} else {
				log.Infof("tag release (%s) matches release ignore constraint '%s': will not update", ctx.Git.Tag, i.String())
				return nil
			}
		}
	}

	if ctx.DryRun {
		log.Info("dry-run: skipping publish")
		return nil
	}

	switch ctx.PublishStrategy {
	case strategies.PublishCommit:
		return publishCommit(ctx)
	case strategies.PublishPullRequest:
		return publishPullRequest(ctx)
	default:
		log.WithFields(log.Fields{
			"strategy": ctx.PublishStrategy,
		}).Error("unsupported publish strategy specified")
		return ErrUnsupportedPublishStrategy
	}
}

func publishCommit(ctx *context.Context) error {
	opts := &client.Options{
		Ref:         ctx.Git.Ref,
		Base:        ctx.Git.Base,
		RepoName:    ctx.Repository.Name,
		RepoOwner:   ctx.Repository.Owner,
		AuthorName:  ctx.Author.Name,
		AuthorEmail: ctx.Author.Email,
	}

	// First, create the remote reference (branch) for the commit, if we are not committing to
	// the master branch.
	if ctx.Git.Ref != "master" {
		if err := ctx.Client.CreateRef(ctx.Context, opts); err != nil {
			return err
		}
	}

	// Update the Chart
	if ctx.Chart.File.HasChanges() {
		if err := ctx.Client.UpdateFile(ctx.Context, opts, ctx.Chart.File.Path, ctx.Release.UpdateCommitMsg, ctx.Chart.File.NewContents); err != nil {
			return err
		}
	} else {
		log.Error("chart has no changes - will not update")
		return ErrNoChartChanges
	}

	for _, f := range ctx.Files {
		if f.HasChanges() {
			if err := ctx.Client.UpdateFile(ctx.Context, opts, f.Path, ctx.Release.UpdateCommitMsg, f.NewContents); err != nil {
				return err
			}
		} else {
			log.WithFields(log.Fields{
				"path": f.Path,
			}).Warn("file has no changes - will not update")
		}
	}
	return nil
}

func publishPullRequest(ctx *context.Context) error {
	if err := publishCommit(ctx); err != nil {
		return err
	}

	opts := &client.Options{
		Ref:         ctx.Git.Ref,
		Base:        ctx.Git.Base,
		RepoName:    ctx.Repository.Name,
		RepoOwner:   ctx.Repository.Owner,
		AuthorName:  ctx.Author.Name,
		AuthorEmail: ctx.Author.Email,
	}

	// Pull request title
	titleTmpl := ctx.Release.PRTitle
	if titleTmpl == "" {
		titleTmpl = templates.DefaultPullRequestTitle
	}

	t, err := template.New("").Parse(titleTmpl)
	if err != nil {
		return err
	}
	buf := bytes.Buffer{}
	err = t.Execute(&buf, ctx)
	if err != nil {
		return err
	}
	title := buf.String()

	// Pull request comment
	commentTmpl := ctx.Release.PRBody
	if commentTmpl == "" {
		commentTmpl = templates.DefaultPullRequestBody
	}

	t, err = template.New("").Parse(commentTmpl)
	if err != nil {
		return err
	}
	buf = bytes.Buffer{}
	err = t.Execute(&buf, ctx)
	if err != nil {
		return err
	}
	body := buf.String()

	log.WithFields(log.Fields{
		"titleTmpl": titleTmpl,
		"bodyTmpl":  commentTmpl,
		"title":     title,
		"body":      body,
	}).Debug("publish: creating pull request")

	return ctx.Client.CreatePullRequest(ctx.Context, opts, title, body)
}
