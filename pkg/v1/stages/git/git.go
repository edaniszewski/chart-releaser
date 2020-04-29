package git

import (
	"github.com/apex/log"
	version "github.com/edaniszewski/chart-releaser/pkg/semver"
	"github.com/edaniszewski/chart-releaser/pkg/utils"
	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
)

// Stage for the "git" step of the update pipeline.
type Stage struct{}

// Name of the stage.
func (Stage) Name() string {
	return "git"
}

// String describes what the stage does.
func (Stage) String() string {
	return "parsing git information"
}

// Run the operations defined for the stage.
func (Stage) Run(ctx *ctx.Context) error {
	log.Debug("looking up git tag")
	tag, err := utils.GetTag()
	if err != nil {
		if err := ctx.CheckDryRun(err); err != nil {
			return err
		}
		tag = "0.0.0"
		log.WithField("tag", tag).Info("using fake tag for dry-run")
	}
	log.WithField("tag", tag).Debug("got git tag")
	ctx.Git.Tag = tag

	// Parse the tag into a Semver. If this fails, we can't continue
	// as we'll be unable to manipulate the chart/app versions correctly.
	v, err := version.Load(tag)
	if err != nil {
		return err
	}

	// The tag version loaded here is the new version for the application
	// being released.
	ctx.App.NewVersion = v

	return nil
}
