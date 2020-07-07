package v1

import (
	"io/ioutil"
	"time"

	"github.com/apex/log"
	"github.com/davecgh/go-spew/spew"
	"github.com/edaniszewski/chart-releaser/pkg/templates"
	v1 "github.com/edaniszewski/chart-releaser/pkg/v1/cfg"
	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/edaniszewski/chart-releaser/pkg/v1/stages/env"
	"github.com/edaniszewski/chart-releaser/pkg/v1/utils"
	"gopkg.in/yaml.v2"
)

// ConfigVersion gets the version of the v1 configuration scheme.
func ConfigVersion() string {
	return v1.ConfigVersion
}

// Updater runs chart updates.
type Updater struct {
	data []byte
	opts UpdateOptions
}

// NewUpdater creates a new Updater.
func NewUpdater(b []byte) *Updater {
	return &Updater{
		data: b,
	}
}

// UpdateOptions provide command line options to the Updater.
type UpdateOptions struct {
	AllowDirty bool
	DryRun     bool
	ShowDiff   bool
	Timeout    time.Duration
}

// AugmentCtx translates the update option values to their corresponding
// fields in a v1 Context.
func (opts *UpdateOptions) AugmentCtx(context *ctx.Context) {
	context.AllowDirty = opts.AllowDirty
	context.DryRun = opts.DryRun
	context.ShowDiff = opts.ShowDiff
}

// Run the v1 chart updater.
func (u *Updater) Run(opts UpdateOptions) error {
	u.opts = opts

	// Load the v1 configuration from the bytes provided to the updater.
	cfg, err := v1.LoadFromBytes(u.data)
	if err != nil {
		return err
	}
	log.Debugf("loaded configuration:\n%v", spew.Sdump(cfg))

	// Create a new context from the loaded configuration.
	context, cancel := ctx.NewWithTimeout(cfg, u.opts.Timeout)
	defer cancel()

	// Set values from the update options onto the Context which will be
	// used throughout the Update process.
	u.opts.AugmentCtx(context)

	if err := cfg.Validate(); err != nil {
		if opts.DryRun {
			log.WithError(err).Warn("dry-run: failed config validation")
		} else {
			return err
		}
	}

	// Run the update pipeline.
	return UpdatePipeline.Run(context)
}

// Formatter runs chart-releaser config formatting.
type Formatter struct {
	data []byte
}

// NewFormatter creates a new Formatter.
func NewFormatter(b []byte) *Formatter {
	return &Formatter{
		data: b,
	}
}

// FormatterOptions provide command line options to the Formatter.
type FormatterOptions struct {
	NoHeader bool
	Path     string
}

// Run the v1 config formatter.
func (f *Formatter) Run(opts FormatterOptions) error {
	cfg, err := v1.LoadFromBytes(f.data)
	if err != nil {
		return err
	}

	yml, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	var data []byte
	if opts.NoHeader {
		data = append(data, yml...)
	} else {
		data = []byte(templates.ConfigHeaderComment)
		data = append(data, []byte("\n")...)
		data = append(data, yml...)
	}
	return ioutil.WriteFile(opts.Path, data, 0644)
}

// Checker runs chart-releaser configuration checks.
type Checker struct {
	data []byte
}

// NewChecker creates a new Checker.
func NewChecker(b []byte) *Checker {
	return &Checker{
		data: b,
	}
}

// CheckerOptions provide command line options to the Checker.
type CheckerOptions struct {
	SkipEnv bool
}

// Run the v1 config checker.
func (c *Checker) Run(opts CheckerOptions) error {
	cfg, err := v1.LoadFromBytes(c.data)
	if err != nil {
		return err
	}

	if !opts.SkipEnv {
		// Evaluate expected environment variables. Note that this does not
		// validate the contents - it just verifies that they are set.
		r, err := utils.ParseRepository(cfg.Chart.Repo)
		if err != nil {
			return err
		}

		context := ctx.New(cfg)
		context.Repository = r

		s := env.Stage{}
		if err := s.Run(context); err != nil {
			return err
		}
	} else {
		log.Debug("skipping check for env vars")
	}

	return cfg.Validate()
}
