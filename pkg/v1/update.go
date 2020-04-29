package v1

import (
	"strings"

	"github.com/apex/log"
	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages/chart"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages/client"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages/config"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages/diff"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages/env"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages/extras"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages/git"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages/publish"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages/render"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages/setup"
	"github.com/fatih/color"
)

// Pipeline defines a list of V1Stages to run in order.
type Pipeline []stages.V1Stage

// Run the stages defined by the Pipeline.
func (p Pipeline) Run(ctx *ctx.Context) error {
	for _, stage := range p {
		log.Info(color.New(color.Bold).Sprintf("%s - %s", strings.ToUpper(stage.Name()), stage.String()))

		err := stage.Run(ctx)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"stage": stage.Name(),
			}).Error("failed running update pipeline stage")
			return err
		}
	}

	if ctx.DryRun {
		ctx.Dump()
		ctx.PrintErrors()
	}
	return nil
}

// UpdatePipeline defines the stages and the order in which to execute them
// in order to perform an update to a Chart.
var UpdatePipeline Pipeline = []stages.V1Stage{
	setup.Stage{},
	config.Stage{},
	env.Stage{},
	git.Stage{},
	client.Stage{},
	chart.Stage{},
	extras.Stage{},
	render.Stage{},
	publish.Stage{},
	diff.Stage{},
}
