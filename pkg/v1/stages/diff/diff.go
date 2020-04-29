package diff

import (
	"fmt"

	"github.com/apex/log"
	"github.com/edaniszewski/chart-releaser/pkg/utils"
	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/mgutz/ansi"
)

// Stage for the "diff" step of the update pipeline.
type Stage struct{}

// Name of the stage.
func (Stage) Name() string {
	return "diff"
}

// String describes what the stage does.
func (Stage) String() string {
	return "displaying changes to chart files"
}

// Run the operations defined for the stage.
func (Stage) Run(ctx *ctx.Context) error {
	if !ctx.ShowDiff {
		log.Info("diff stage not enabled - skipping")
		return nil
	}

	_, _ = fmt.Fprint(ctx.Out, ansi.Color("+----start diff\n", "white"))

	utils.PrintDiff(
		ctx.Out,
		"Chart.yaml",
		string(ctx.Chart.File.PreviousContents),
		string(ctx.Chart.File.NewContents),
	)

	for _, f := range ctx.Files {
		utils.PrintDiff(
			ctx.Out,
			f.Path,
			string(f.PreviousContents),
			string(f.NewContents),
		)
	}

	_, _ = fmt.Fprint(ctx.Out, ansi.Color("+----end diff\n", "white"))
	return nil
}
