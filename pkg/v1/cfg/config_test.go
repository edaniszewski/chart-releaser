package v1

import (
	"testing"

	"github.com/edaniszewski/chart-releaser/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func TestLoadFromBytes(t *testing.T) {
	b := []byte(`
version: v1
chart:
  name: test-chart
`)

	c, err := LoadFromBytes(b)
	assert.NoError(t, err)
	assert.Equal(t, c.Version, "v1")
	assert.Equal(t, c.Chart.Name, "test-chart")
}

func TestLoadFromBytes_Error(t *testing.T) {
	_, err := LoadFromBytes([]byte{0x00, 0x01})
	assert.EqualError(t, err, "error converting YAML to JSON: yaml: control characters are not allowed")
}

func TestConfig_GetVersion(t *testing.T) {
	c := Config{
		Version: "v1",
	}
	assert.Equal(t, c.GetVersion(), "v1")
}

func TestConfig_GetVersion_Empty(t *testing.T) {
	c := Config{}
	assert.Equal(t, c.GetVersion(), "")
}

func TestConfig_Validate(t *testing.T) {
	cfg := Config{
		Version: "v1",
		Chart: &ChartConfig{
			Name: "test-chart",
			Repo: "github.com/test-chart",
		},
		Publish: &PublishConfig{
			PR: &PublishPRConfig{},
		},
		Commit: &CommitConfig{
			Author: &CommitAuthorConfig{
				Name:  "test-user",
				Email: "test-email",
			},
			Templates: &CommitTemplateConfig{},
		},
		Release: &ReleaseConfig{
			Strategy: "default",
		},
		Extras: []*ExtrasConfig{
			{
				Path: "test path",
				Updates: []*SearchReplace{
					{
						Search:  "old value",
						Replace: "new value",
					},
				},
			},
		},
	}

	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestConfig_Validate_Error(t *testing.T) {
	cfg := Config{
		Version: "v0",
		Chart:   &ChartConfig{},
		Publish: &PublishConfig{
			PR:     &PublishPRConfig{},
			Commit: &PublishCommitConfig{},
		},
		Commit: &CommitConfig{
			Author: &CommitAuthorConfig{
				Name: "test-user",
			},
			Templates: &CommitTemplateConfig{},
		},
		Release: &ReleaseConfig{
			Strategy: "invalid",
		},
		Extras: []*ExtrasConfig{
			{
				Path: "test path",
			},
		},
	}

	err := cfg.Validate()
	assert.Error(t, err)

	collector, ok := err.(*errs.Collector)
	assert.True(t, ok, "error is an instance of errs.Collector")
	assert.Equal(t, 7, collector.Count())
	assert.EqualError(t, err, `
Errors:
 • error: using v1 config parser for non v1 config
 • required option 'chart.name' missing from config
 • required option 'chart.repo' missing from config
 • invalid publish config: cannot define both 'commit' and 'pr' blocks
 • commit author specifies name, but no email
 • invalid release strategy 'invalid', should be one of: [major minor patch default]
 • extras config specifies path but no options for search/replace updates

`)
}

func TestChartConfig_validate(t *testing.T) {
	cfg := ChartConfig{
		Name: "test-chart",
		Repo: "github.com/test/repo",
	}

	err := cfg.validate()
	assert.NoError(t, err)
}

func TestChartConfig_validateErrors(t *testing.T) {
	cfg := ChartConfig{}

	err := cfg.validate()
	assert.Error(t, err)

	collector, ok := err.(*errs.Collector)
	assert.True(t, ok, "error is an instance of errs.Collector")
	assert.Equal(t, 2, collector.Count())
	assert.EqualError(t, err, `
Errors:
 • required option 'chart.name' missing from config
 • required option 'chart.repo' missing from config

`)
}

func TestChartConfig_validateErrors2(t *testing.T) {
	cfg := ChartConfig{
		Repo: "some-repo",
	}

	err := cfg.validate()
	assert.Error(t, err)

	collector, ok := err.(*errs.Collector)
	assert.True(t, ok, "error is an instance of errs.Collector")
	assert.Equal(t, 2, collector.Count())
	assert.EqualError(t, err, `
Errors:
 • required option 'chart.name' missing from config
 • unsupported repo specified in 'chart.repo'. currently supported repos include: github.com

`)
}

func TestPublishConfig_validate(t *testing.T) {
	cfg := PublishConfig{
		PR: &PublishPRConfig{},
	}

	err := cfg.validate()
	assert.NoError(t, err)
}

func TestPublishConfig_validateErrors(t *testing.T) {
	cfg := PublishConfig{
		PR:     &PublishPRConfig{},
		Commit: &PublishCommitConfig{},
	}

	err := cfg.validate()
	assert.Error(t, err)

	collector, ok := err.(*errs.Collector)
	assert.True(t, ok, "error is an instance of errs.Collector")
	assert.Equal(t, 1, collector.Count())
	assert.EqualError(t, err, `
Errors:
 • invalid publish config: cannot define both 'commit' and 'pr' blocks

`)
}

func TestCommitConfig_validate(t *testing.T) {
	cfg := CommitConfig{
		Author: &CommitAuthorConfig{
			Name:  "test-user",
			Email: "test-email",
		},
		Templates: &CommitTemplateConfig{},
	}

	err := cfg.validate()
	assert.NoError(t, err)
}

func TestCommitConfig_validateErrors(t *testing.T) {
	cfg := CommitConfig{
		Author: &CommitAuthorConfig{
			Name: "test-user",
		},
		Templates: &CommitTemplateConfig{},
	}

	err := cfg.validate()
	assert.Error(t, err)

	collector, ok := err.(*errs.Collector)
	assert.True(t, ok, "error is an instance of errs.Collector")
	assert.Equal(t, 1, collector.Count())
	assert.EqualError(t, err, `
Errors:
 • commit author specifies name, but no email

`)
}

func TestCommitTemplateConfig_validate(t *testing.T) {
	cfg := CommitTemplateConfig{}

	err := cfg.validate()
	assert.NoError(t, err)
}

func TestCommitTemplateConfig_validateErrors(t *testing.T) {
	// TODO (etd): Nothing currently validated.
}

func TestReleaseConfig_validate(t *testing.T) {
	cfg := ReleaseConfig{
		Strategy: "default",
	}

	err := cfg.validate()
	assert.NoError(t, err)
}

func TestReleaseConfig_validateErrors(t *testing.T) {
	cfg := ReleaseConfig{
		Strategy: "invalid",
	}

	err := cfg.validate()
	assert.Error(t, err)

	collector, ok := err.(*errs.Collector)
	assert.True(t, ok, "error is an instance of errs.Collector")
	assert.Equal(t, 1, collector.Count())
	assert.EqualError(t, err, `
Errors:
 • invalid release strategy 'invalid', should be one of: [major minor patch default]

`)
}

func TestPublishCommitConfig_validate(t *testing.T) {
	cfg := PublishCommitConfig{}

	err := cfg.validate()
	assert.NoError(t, err)
}

func TestPublishCommitConfig_validateErrors(t *testing.T) {
	// TODO (etd): Nothing currently validated.
}

func TestPublishPRConfig_validate(t *testing.T) {
	cfg := PublishPRConfig{}

	err := cfg.validate()
	assert.NoError(t, err)
}

func TestPublishPRConfig_validateErrors(t *testing.T) {
	// TODO (etd): Nothing currently validated.
}

func TestCommitAuthorConfig_validate(t *testing.T) {
	cfg := CommitAuthorConfig{
		Name:  "test-user",
		Email: "test-email",
	}

	err := cfg.validate()
	assert.NoError(t, err)
}

func TestCommitAuthorConfig_validateErrors(t *testing.T) {
	cfg := CommitAuthorConfig{
		Name: "test-user",
	}

	err := cfg.validate()
	assert.Error(t, err)

	collector, ok := err.(*errs.Collector)
	assert.True(t, ok, "error is an instance of errs.Collector")
	assert.Equal(t, 1, collector.Count())
	assert.EqualError(t, err, `
Errors:
 • commit author specifies name, but no email

`)
}

func TestCommitAuthorConfig_validateErrors2(t *testing.T) {
	cfg := CommitAuthorConfig{
		Email: "test-email",
	}

	err := cfg.validate()
	assert.Error(t, err)

	collector, ok := err.(*errs.Collector)
	assert.True(t, ok, "error is an instance of errs.Collector")
	assert.Equal(t, 1, collector.Count())
	assert.EqualError(t, err, `
Errors:
 • commit author specified email, but no name

`)
}

func TestExtrasConfig_validate(t *testing.T) {
	cfg := ExtrasConfig{
		Path: "test",
		Updates: []*SearchReplace{
			{
				Search:  "old value",
				Replace: "new value",
			},
		},
	}

	err := cfg.validate()
	assert.NoError(t, err)
}

func TestExtrasConfig_validateErrors(t *testing.T) {
	cfg := ExtrasConfig{
		Path: "test",
	}

	err := cfg.validate()
	assert.Error(t, err)

	collector, ok := err.(*errs.Collector)
	assert.True(t, ok, "error is an instance of errs.Collector")
	assert.Equal(t, 1, collector.Count())
	assert.EqualError(t, err, `
Errors:
 • extras config specifies path but no options for search/replace updates

`)
}

func TestExtrasConfig_validateErrors2(t *testing.T) {
	cfg := ExtrasConfig{
		Path: "test",
		Updates: []*SearchReplace{
			{Limit: -1},
		},
	}

	err := cfg.validate()
	assert.Error(t, err)

	collector, ok := err.(*errs.Collector)
	assert.True(t, ok, "error is an instance of errs.Collector")
	assert.Equal(t, 3, collector.Count())
	assert.EqualError(t, err, `
Errors:
 • search and replace config has no search string set
 • search and replace config has no replace string set
 • search and replace config must have a non-negative limit

`)
}

func TestSearchReplaceConfig_validate(t *testing.T) {
	cfg := SearchReplace{
		Search:  "old value",
		Replace: "new value",
	}

	err := cfg.validate()
	assert.NoError(t, err)
}

func TestSearchReplaceConfig_validateErrors(t *testing.T) {
	cfg := SearchReplace{
		Limit: -1,
	}

	err := cfg.validate()
	assert.Error(t, err)

	collector, ok := err.(*errs.Collector)
	assert.True(t, ok, "error is an instance of errs.Collector")
	assert.Equal(t, 3, collector.Count())
	assert.EqualError(t, err, `
Errors:
 • search and replace config has no search string set
 • search and replace config has no replace string set
 • search and replace config must have a non-negative limit

`)
}
