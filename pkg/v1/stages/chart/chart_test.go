package chart

import (
	"errors"
	"math"
	"testing"

	"github.com/edaniszewski/chart-releaser/internal/testutils"
	"github.com/edaniszewski/chart-releaser/pkg/strategies"
	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/stretchr/testify/assert"
)

func TestStage_Name(t *testing.T) {
	assert.Equal(t, "chart", Stage{}.Name())
}

func TestStage_String(t *testing.T) {
	assert.Equal(t, "updating helm chart", Stage{}.String())
}

func TestStage_Run(t *testing.T) {
	context := ctx.Context{
		Chart: ctx.Chart{
			Name:    "test-chart",
			SubPath: "charts",
		},
		App: ctx.App{
			NewVersion: testutils.NewSemver(t, "0.3.0"),
		},
		UpdateStrategy: strategies.UpdateDefault,
		Client: &testutils.FakeClient{
			FileData: `
apiVersion: v1
name: test-chart
version: 0.1.2
appVersion: 0.2.3
`,
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)

	// Context relating to Chart data
	assert.Equal(t, "test-chart", context.Chart.Name)
	assert.Equal(t, "charts", context.Chart.SubPath)
	assert.Equal(t, "0.1.3", context.Chart.NewVersion.String())
	assert.Equal(t, "0.1.2", context.Chart.PreviousVersion.String())
	assert.Equal(t, "charts/Chart.yaml", context.Chart.File.Path)
	assert.Equal(t, "apiVersion: v1\nappVersion: 0.3.0\nname: test-chart\nversion: 0.1.3\n", string(context.Chart.File.NewContents))
	assert.Equal(t, "apiVersion: v1\nappVersion: 0.2.3\nname: test-chart\nversion: 0.1.2\n", string(context.Chart.File.PreviousContents))

	assert.Equal(t, "0.3.0", context.App.NewVersion.String())
	assert.Equal(t, "0.2.3", context.App.PreviousVersion.String())

	// Other context info. We don't expect this to change.
	assert.Equal(t, "", context.Author.Name)
	assert.Equal(t, "", context.Author.Email)
	assert.Equal(t, "", context.Git.Base)
	assert.Equal(t, "", context.Git.Ref)
	assert.Equal(t, "", context.Git.Tag)
	assert.Equal(t, "", context.Release.UpdateCommitMsg)
	assert.Equal(t, "", context.Release.PRBody)
	assert.Equal(t, "", context.Release.PRTitle)
	assert.Equal(t, "", context.Repository.Name)
	assert.Equal(t, "", context.Repository.Owner)
	assert.Equal(t, ctx.RepoType(""), context.Repository.Type)
	assert.Len(t, context.Files, 0)
}

func TestStage_RunChartGetError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{
			GetFileError: []error{
				errors.New("testing"),
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "testing")

	// Context relating to Chart data
	assert.Equal(t, "", context.Chart.Name)
	assert.Equal(t, "", context.Chart.SubPath)
	assert.Equal(t, "0.0.0", context.Chart.NewVersion.String())
	assert.Equal(t, "0.0.0", context.Chart.PreviousVersion.String())
	assert.Equal(t, "Chart.yaml", context.Chart.File.Path)
	assert.Equal(t, "", string(context.Chart.File.NewContents))
	assert.Equal(t, "", string(context.Chart.File.PreviousContents))

	assert.Equal(t, "0.0.0", context.App.NewVersion.String())
	assert.Equal(t, "0.0.0", context.App.PreviousVersion.String())
}

func TestStage_RunUnmarshalChartError(t *testing.T) {
	context := ctx.Context{
		Chart: ctx.Chart{
			Name:    "test-chart",
			SubPath: "charts",
		},
		App: ctx.App{
			NewVersion: testutils.NewSemver(t, "0.3.0"),
		},
		UpdateStrategy: strategies.UpdateDefault,
		Client: &testutils.FakeClient{
			FileData: `
apiVersion: [
`,
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "error converting YAML to JSON: yaml: line 2: did not find expected node content")

	// Context relating to Chart data
	assert.Equal(t, "test-chart", context.Chart.Name)
	assert.Equal(t, "charts", context.Chart.SubPath)
	assert.Equal(t, "0.0.0", context.Chart.NewVersion.String())
	assert.Equal(t, "0.0.0", context.Chart.PreviousVersion.String())
	assert.Equal(t, "charts/Chart.yaml", context.Chart.File.Path)
	assert.Equal(t, "", string(context.Chart.File.NewContents))
	assert.Equal(t, "", string(context.Chart.File.PreviousContents))

	assert.Equal(t, "0.3.0", context.App.NewVersion.String())
	assert.Equal(t, "0.0.0", context.App.PreviousVersion.String())
}

func TestStage_RunNoChartVersionError(t *testing.T) {
	context := ctx.Context{
		Chart: ctx.Chart{
			Name:    "test-chart",
			SubPath: "charts",
		},
		App: ctx.App{
			NewVersion: testutils.NewSemver(t, "0.3.0"),
		},
		UpdateStrategy: strategies.UpdateDefault,
		Client: &testutils.FakeClient{
			FileData: `
apiVersion: v1
name: test-chart
appVersion: 0.2.3
`,
		},
	}

	err := Stage{}.Run(&context)
	assert.Equal(t, ErrNoChartVersion, err)

	// Context relating to Chart data
	assert.Equal(t, "test-chart", context.Chart.Name)
	assert.Equal(t, "charts", context.Chart.SubPath)
	assert.Equal(t, "0.0.0", context.Chart.NewVersion.String())
	assert.Equal(t, "0.0.0", context.Chart.PreviousVersion.String())
	assert.Equal(t, "charts/Chart.yaml", context.Chart.File.Path)
	assert.Equal(t, "", string(context.Chart.File.NewContents))
	assert.Equal(t, "apiVersion: v1\nappVersion: 0.2.3\nname: test-chart\n", string(context.Chart.File.PreviousContents))

	assert.Equal(t, "0.3.0", context.App.NewVersion.String())
	assert.Equal(t, "0.0.0", context.App.PreviousVersion.String())
}

func TestStage_RunNoAppVersionError(t *testing.T) {
	context := ctx.Context{
		Chart: ctx.Chart{
			Name:    "test-chart",
			SubPath: "charts",
		},
		App: ctx.App{
			NewVersion: testutils.NewSemver(t, "0.3.0"),
		},
		UpdateStrategy: strategies.UpdateDefault,
		Client: &testutils.FakeClient{
			FileData: `
apiVersion: v1
name: test-chart
version: 0.1.2
`,
		},
	}

	err := Stage{}.Run(&context)
	assert.Equal(t, err, ErrNoAppVersion)

	// Context relating to Chart data
	assert.Equal(t, "test-chart", context.Chart.Name)
	assert.Equal(t, "charts", context.Chart.SubPath)
	assert.Equal(t, "0.0.0", context.Chart.NewVersion.String())
	assert.Equal(t, "0.0.0", context.Chart.PreviousVersion.String())
	assert.Equal(t, "charts/Chart.yaml", context.Chart.File.Path)
	assert.Equal(t, "", string(context.Chart.File.NewContents))
	assert.Equal(t, "apiVersion: v1\nname: test-chart\nversion: 0.1.2\n", string(context.Chart.File.PreviousContents))

	assert.Equal(t, "0.3.0", context.App.NewVersion.String())
	assert.Equal(t, "0.0.0", context.App.PreviousVersion.String())
}

func TestStage_RunBadChartVersionError(t *testing.T) {
	context := ctx.Context{
		Chart: ctx.Chart{
			Name:    "test-chart",
			SubPath: "charts",
		},
		App: ctx.App{
			NewVersion: testutils.NewSemver(t, "0.3.0"),
		},
		UpdateStrategy: strategies.UpdateDefault,
		Client: &testutils.FakeClient{
			FileData: `
apiVersion: v1
name: test-chart
version: not-a-semver
appVersion: 0.2.3
`,
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "No Major.Minor.Patch elements found")

	// Context relating to Chart data
	assert.Equal(t, "test-chart", context.Chart.Name)
	assert.Equal(t, "charts", context.Chart.SubPath)
	assert.Equal(t, "0.0.0", context.Chart.NewVersion.String())
	assert.Equal(t, "0.0.0", context.Chart.PreviousVersion.String())
	assert.Equal(t, "charts/Chart.yaml", context.Chart.File.Path)
	assert.Equal(t, "", string(context.Chart.File.NewContents))
	assert.Equal(t, "apiVersion: v1\nappVersion: 0.2.3\nname: test-chart\nversion: not-a-semver\n", string(context.Chart.File.PreviousContents))

	assert.Equal(t, "0.3.0", context.App.NewVersion.String())
	assert.Equal(t, "0.0.0", context.App.PreviousVersion.String())
}

func TestStage_RunBadAppVersionError(t *testing.T) {
	context := ctx.Context{
		Chart: ctx.Chart{
			Name:    "test-chart",
			SubPath: "charts",
		},
		App: ctx.App{
			NewVersion: testutils.NewSemver(t, "0.3.0"),
		},
		UpdateStrategy: strategies.UpdateDefault,
		Client: &testutils.FakeClient{
			FileData: `
apiVersion: v1
name: test-chart
version: 0.1.2
appVersion: not-a-semver
`,
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "No Major.Minor.Patch elements found")

	// Context relating to Chart data
	assert.Equal(t, "test-chart", context.Chart.Name)
	assert.Equal(t, "charts", context.Chart.SubPath)
	assert.Equal(t, "0.0.0", context.Chart.NewVersion.String())
	assert.Equal(t, "0.1.2", context.Chart.PreviousVersion.String())
	assert.Equal(t, "charts/Chart.yaml", context.Chart.File.Path)
	assert.Equal(t, "", string(context.Chart.File.NewContents))
	assert.Equal(t, "apiVersion: v1\nappVersion: not-a-semver\nname: test-chart\nversion: 0.1.2\n", string(context.Chart.File.PreviousContents))

	assert.Equal(t, "0.3.0", context.App.NewVersion.String())
	assert.Equal(t, "0.0.0", context.App.PreviousVersion.String())
}

func TestStage_RunUpdateReleaseError(t *testing.T) {
	context := ctx.Context{
		Chart: ctx.Chart{
			Name:    "test-chart",
			SubPath: "charts",
		},
		App: ctx.App{
			NewVersion: testutils.NewSemver(t, "0.2.0"), // new version less than previous
		},
		UpdateStrategy: strategies.UpdateDefault,
		Client: &testutils.FakeClient{
			FileData: `
apiVersion: v1
name: test-chart
version: 0.1.2
appVersion: 0.2.3
`,
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "version is less than the starting base version")

	// Context relating to Chart data
	assert.Equal(t, "test-chart", context.Chart.Name)
	assert.Equal(t, "charts", context.Chart.SubPath)
	assert.Equal(t, "0.0.0", context.Chart.NewVersion.String())
	assert.Equal(t, "0.1.2", context.Chart.PreviousVersion.String())
	assert.Equal(t, "charts/Chart.yaml", context.Chart.File.Path)
	assert.Equal(t, "", string(context.Chart.File.NewContents))
	assert.Equal(t, "apiVersion: v1\nappVersion: 0.2.3\nname: test-chart\nversion: 0.1.2\n", string(context.Chart.File.PreviousContents))

	assert.Equal(t, "0.2.0", context.App.NewVersion.String())
	assert.Equal(t, "0.2.3", context.App.PreviousVersion.String())
}

func TestStage_RunDryRunErrors(t *testing.T) {
	context := ctx.Context{
		DryRun: true,
		Chart: ctx.Chart{
			Name:    "test-chart",
			SubPath: "charts",
		},
		App: ctx.App{
			NewVersion: testutils.NewSemver(t, "0.3.0"),
		},
		UpdateStrategy: strategies.UpdateDefault,
		Client: &testutils.FakeClient{
			FileData: `
apiVersion: v1
name: test-chart
`,
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
	assert.EqualError(t, context.Errors(), "\nErrors:\n • chart does not specify a version\n • chart does not specify an appVersion\n\n")
}

func TestStage_RunDryRunErrors2(t *testing.T) {
	context := ctx.Context{
		DryRun: true,
		Chart: ctx.Chart{
			Name:    "test-chart",
			SubPath: "charts",
		},
		App: ctx.App{
			NewVersion: testutils.NewSemver(t, "0.2.0"),
		},
		UpdateStrategy: strategies.UpdateDefault,
		Client: &testutils.FakeClient{
			FileData: `
apiVersion: v1
name: test-chart
version: 0.1.2
appVersion: 0.2.3
`,
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)
	assert.EqualError(t, context.Errors(), "\nErrors:\n • version is less than the starting base version\n\n")
}

func TestMarshalContents(t *testing.T) {
	actual, err := marshalContents(&struct {
		Foo string
		Bar int
	}{
		Foo: "test",
		Bar: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, "Bar: 1\nFoo: test\n", string(actual))
}

func TestMarshalContentsError(t *testing.T) {
	_, err := marshalContents(math.NaN())
	assert.EqualError(t, err, "json: unsupported value: NaN")
}
