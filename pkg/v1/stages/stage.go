package stages

import (
	"fmt"

	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
)

// V1Stage defines an interface for all stages of a pipeline
// which operate on a v1 Context.
type V1Stage interface {
	fmt.Stringer

	Name() string
	Run(ctx *ctx.Context) error
}
