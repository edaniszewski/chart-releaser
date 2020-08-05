package chart

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/edaniszewski/chart-releaser/pkg/client"
	version "github.com/edaniszewski/chart-releaser/pkg/semver"
	"github.com/edaniszewski/chart-releaser/pkg/strategies"
	context "github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
)

// Errors for the chart stage.
var (
	ErrNoChartVersion = errors.New("chart does not specify a version")
	ErrNoAppVersion   = errors.New("chart does not specify an appVersion")
)

// Stage for the "chart" step of the update pipeline.
type Stage struct{}

// Name of the stage.
func (Stage) Name() string {
	return "chart"
}

// String describes what the stage does.
func (Stage) String() string {
	return "updating helm chart"
}

// Run the operations defined for the stage.
func (Stage) Run(ctx *context.Context) error {
	// Check that the given path is for a Chart.yaml file. If not, consider it
	// a directory and append "Chart.yaml" to the path.
	path := ctx.Chart.SubPath
	if !strings.HasSuffix(path, "Chart.yaml") && !strings.HasSuffix(path, "Chart.yml") {
		path = filepath.Join(path, "Chart.yaml")
	}

	ctx.Chart.File = context.File{
		Path: path,
	}

	opts := &client.Options{
		RepoName:  ctx.Repository.Name,
		RepoOwner: ctx.Repository.Owner,
	}

	// Get the chart.
	raw, err := ctx.Client.GetFile(ctx.Context, opts, path)
	if err != nil {
		return err
	}

	chartMeta, err := chartutil.UnmarshalChartfile([]byte(raw))
	if err != nil {
		return err
	}

	ctx.Chart.File.PreviousContents, err = marshalContents(chartMeta)
	if err != nil {
		return err
	}

	// Get the chart version and the app version defined in the Chart
	if chartMeta.Version == "" {
		if err := ctx.CheckDryRun(ErrNoChartVersion); err != nil {
			return err
		}
		chartMeta.Version = "0.0.0"
		log.WithField("version", chartMeta.Version).Warn("dry-run: using placeholder for chart version")
	}
	if chartMeta.AppVersion == "" {
		if err := ctx.CheckDryRun(ErrNoAppVersion); err != nil {
			return err
		}
		chartMeta.AppVersion = "0.0.0"
		log.WithField("appVersion", chartMeta.AppVersion).Warn("dry-run: using placeholder for appVersion")
	}

	ctx.Chart.PreviousVersion, err = version.Load(chartMeta.Version)
	if err != nil {
		return err
	}

	ctx.App.PreviousVersion, err = version.Load(chartMeta.AppVersion)
	if err != nil {
		return err
	}

	// Determine the new version of the chart.
	ctx.Chart.NewVersion, err = strategies.UpdateRelease(&strategies.UpdateCtx{
		OldAppVersion:   &ctx.App.PreviousVersion,
		NewAppVersion:   &ctx.App.NewVersion,
		OldChartVersion: &ctx.Chart.PreviousVersion,
		Strategy:        ctx.UpdateStrategy,
	})
	if err != nil {
		if err := ctx.CheckDryRun(err); err != nil {
			return err
		}
		v, _ := version.Load("0.1.0")
		ctx.Chart.NewVersion = v
		log.WithField("version", "0.1.0").Warn("dry-run: using placeholder for new chart version")
	}

	// Update the chart metadata struct with the new values.
	chartMeta.Version = ctx.Chart.NewVersion.String()
	chartMeta.AppVersion = ctx.App.NewVersion.String()

	// Serialize the chart to yaml bytes for writing.
	ctx.Chart.File.NewContents, err = marshalContents(chartMeta)
	if err != nil {
		return err
	}
	return nil
}

// marshalContents marshals the Chart metadata first to JSON then to YAML. This is a
// bit of a hack because the struct does not include annotations for YAML, do direct
// marshalling to YAML will result in undesired fields, and inconsistent field names.
func marshalContents(v interface{}) ([]byte, error) {
	c, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return yaml.JSONToYAML(c)
}
