package config

import (
	"testing"

	"github.com/edaniszewski/chart-releaser/internal/testutils"
	"github.com/edaniszewski/chart-releaser/pkg/strategies"
	"github.com/edaniszewski/chart-releaser/pkg/templates"
	v1 "github.com/edaniszewski/chart-releaser/pkg/v1/cfg"
	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/stretchr/testify/assert"
)

func TestStage_Name(t *testing.T) {
	assert.Equal(t, "config", Stage{}.Name())
}

func TestStage_String(t *testing.T) {
	assert.Equal(t, "loading context configuration", Stage{}.String())
}

func TestStage_Run(t *testing.T) {
	context := ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{
				Strategy: "minor",
			},
			Commit: &v1.CommitConfig{
				Author: &v1.CommitAuthorConfig{
					Name:  "test-user",
					Email: "test-email",
				},
			},
			Chart: &v1.ChartConfig{
				Name: "test-chart",
				Path: "test-path",
				Repo: "github.com/edaniszewski/charts-test",
			},
			Publish: &v1.PublishConfig{
				PR: &v1.PublishPRConfig{},
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)

	// Chart
	assert.Equal(t, "test-chart", context.Chart.Name)
	assert.True(t, context.Chart.NewVersion.Equals(testutils.NewSemverP(t, "0.0.0")))
	assert.True(t, context.Chart.PreviousVersion.Equals(testutils.NewSemverP(t, "0.0.0")))
	assert.Equal(t, "test-path", context.Chart.SubPath)
	assert.Equal(t, ctx.File{
		Path:             "",
		PreviousContents: []byte(nil),
		NewContents:      []byte(nil),
	}, context.Chart.File)

	// Author
	assert.Equal(t, "test-user", context.Author.Name)
	assert.Equal(t, "test-email", context.Author.Email)

	// Repository
	assert.Equal(t, "charts-test", context.Repository.Name)
	assert.Equal(t, "edaniszewski", context.Repository.Owner)
	assert.Equal(t, ctx.RepoGithub, context.Repository.Type)

	// App
	assert.True(t, context.App.PreviousVersion.Equals(testutils.NewSemverP(t, "0.0.0")))
	assert.True(t, context.App.NewVersion.Equals(testutils.NewSemverP(t, "0.0.0")))

	// Files
	assert.Len(t, context.Files, 0)

	// Git
	assert.Equal(t, "master", context.Git.Base)
	assert.Equal(t, "chartreleaser/{{ .Chart.Name }}/{{ .Chart.NewVersion }}", context.Git.Ref)
	assert.Equal(t, "", context.Git.Tag)

	// Release
	assert.Equal(t, "Bump {{ .Chart.Name }} Chart from {{ .Chart.PreviousVersion }} to {{ .Chart.NewVersion }}", context.Release.PRTitle)
	assert.Equal(t, "Bumps the {{ .Chart.Name }} Helm Chart from {{ .Chart.PreviousVersion }} to {{ .Chart.NewVersion }}.\n\n{{ if .Files }}The following files have also been updated:\n{{ range .Files }}- <pre>{{ .Path }}</pre>\n{{ end }}{{ end }}\n---\n*This PR was opened using [chart-releaser](https://github.com/edaniszewski/chart-releaser)*\n", context.Release.PRBody)
	assert.Equal(t, "[{{ .Chart.Name }}] bump chart to {{ .Chart.NewVersion }} for new application release ({{ .App.NewVersion }})", context.Release.UpdateCommitMsg)

	// Other
	assert.Equal(t, "", context.Token)
	assert.Equal(t, strategies.PublishPullRequest, context.PublishStrategy)
	assert.Equal(t, strategies.UpdateMinor, context.UpdateStrategy)
	assert.Equal(t, false, context.AllowDirty)
	assert.Equal(t, false, context.DryRun)
	assert.Equal(t, false, context.ShowDiff)
	assert.Nil(t, context.Client)
}

func TestStage_Run_ErrorUpdateStrategy(t *testing.T) {
	context := ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{
				Strategy: "unsupported",
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "not a valid strategy: unsupported")
}

func TestStage_Run_ErrorBadRepository(t *testing.T) {
	context := ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{
				Strategy: "default",
			},
			Commit: &v1.CommitConfig{
				Author: &v1.CommitAuthorConfig{
					Name:  "test-user",
					Email: "test-email",
				},
			},
			Chart: &v1.ChartConfig{
				Name: "test-chart",
				Path: "test-path",
				Repo: "a/b/c/d/e",
			},
			Publish: &v1.PublishConfig{
				PR: &v1.PublishPRConfig{},
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.EqualError(t, err, "unexpected repository string format - should be in the form of REPO/OWNER/NAME")
}

func TestLoadUpdateStrategy(t *testing.T) {
	tests := []struct {
		name     string
		strategy string
		expected strategies.UpdateStrategy
	}{
		{
			name:     "default strategy",
			strategy: "",
			expected: strategies.UpdateDefault,
		},
		{
			name:     "strategy: major",
			strategy: "major",
			expected: strategies.UpdateMajor,
		},
		{
			name:     "strategy: minor",
			strategy: "minor",
			expected: strategies.UpdateMinor,
		},
		{
			name:     "strategy: patch",
			strategy: "patch",
			expected: strategies.UpdatePatch,
		},
		{
			name:     "strategy: default",
			strategy: "default",
			expected: strategies.UpdateDefault,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context := ctx.Context{
				Config: &v1.Config{
					Release: &v1.ReleaseConfig{
						Strategy: test.strategy,
					},
				},
			}
			assert.Empty(t, context.UpdateStrategy)

			err := loadUpdateStrategy(&context)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, context.UpdateStrategy)
		})
	}
}

func TestLoadUpdateStrategy_Error(t *testing.T) {
	err := loadUpdateStrategy(&ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{
				Strategy: "invalid",
			},
		},
	})
	assert.EqualError(t, err, "not a valid strategy: invalid")
}

func TestLoadPublishStrategy(t *testing.T) {
	tests := []struct {
		name     string
		cfg      v1.PublishConfig
		expected strategies.PublishStrategy
	}{
		{
			name:     "default strategy",
			cfg:      v1.PublishConfig{},
			expected: strategies.PublishPullRequest,
		},
		{
			name: "strategy: commit",
			cfg: v1.PublishConfig{
				Commit: &v1.PublishCommitConfig{},
			},
			expected: strategies.PublishCommit,
		},
		{
			name: "strategy: pull request",
			cfg: v1.PublishConfig{
				PR: &v1.PublishPRConfig{},
			},
			expected: strategies.PublishPullRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context := ctx.Context{
				Config: &v1.Config{
					Publish: &test.cfg,
				},
			}
			assert.Empty(t, context.PublishStrategy)

			err := loadPublishStrategy(&context)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, context.PublishStrategy)
		})
	}
}

func TestLoadTemplateStringsDefaultsForCommit(t *testing.T) {
	context := &ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{},
			Commit:  &v1.CommitConfig{},
			Publish: &v1.PublishConfig{
				Commit: &v1.PublishCommitConfig{},
			},
		},
		PublishStrategy: strategies.PublishCommit,
	}

	err := loadTemplateStrings(context)
	assert.NoError(t, err)

	assert.Equal(t, templates.DefaultUpdateCommitMessage, context.Release.UpdateCommitMsg)
	assert.Equal(t, "", context.Release.PRTitle)
	assert.Equal(t, "", context.Release.PRBody)
	assert.Equal(t, "master", context.Git.Ref)
	assert.Equal(t, "master", context.Git.Base)
}

func TestLoadTemplateStringsCustomForCommit(t *testing.T) {
	context := &ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{},
			Commit: &v1.CommitConfig{
				Templates: &v1.CommitTemplateConfig{
					Update: "test update",
				},
			},
			Publish: &v1.PublishConfig{
				Commit: &v1.PublishCommitConfig{
					Branch: "test-branch",
					Base:   "test-base",
				},
			},
		},
		PublishStrategy: strategies.PublishCommit,
	}

	err := loadTemplateStrings(context)
	assert.NoError(t, err)

	assert.Equal(t, "test update", context.Release.UpdateCommitMsg)
	assert.Equal(t, "", context.Release.PRTitle)
	assert.Equal(t, "", context.Release.PRBody)
	assert.Equal(t, "test-branch", context.Git.Ref)
	assert.Equal(t, "test-base", context.Git.Base)
}

func TestLoadTemplateStringsDefaultsForPR(t *testing.T) {
	context := &ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{},
			Commit:  &v1.CommitConfig{},
			Publish: &v1.PublishConfig{
				PR: &v1.PublishPRConfig{},
			},
		},
		PublishStrategy: strategies.PublishPullRequest,
	}

	err := loadTemplateStrings(context)
	assert.NoError(t, err)

	assert.Equal(t, templates.DefaultUpdateCommitMessage, context.Release.UpdateCommitMsg)
	assert.Equal(t, templates.DefaultPullRequestTitle, context.Release.PRTitle)
	assert.Equal(t, templates.DefaultPullRequestBody, context.Release.PRBody)
	assert.Equal(t, templates.DefaultBranchName, context.Git.Ref)
	assert.Equal(t, "master", context.Git.Base)
}

func TestLoadTemplateStringsCustomForPR(t *testing.T) {
	context := &ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{},
			Commit: &v1.CommitConfig{
				Templates: &v1.CommitTemplateConfig{
					Update: "test update",
				},
			},
			Publish: &v1.PublishConfig{
				PR: &v1.PublishPRConfig{
					BranchTemplate: "branch-template",
					Base:           "test-base",
					TitleTemplate:  "title-template",
					BodyTemplate:   "body-template",
				},
			},
		},
		PublishStrategy: strategies.PublishPullRequest,
	}

	err := loadTemplateStrings(context)
	assert.NoError(t, err)

	assert.Equal(t, "test update", context.Release.UpdateCommitMsg)
	assert.Equal(t, "title-template", context.Release.PRTitle)
	assert.Equal(t, "body-template", context.Release.PRBody)
	assert.Equal(t, "branch-template", context.Git.Ref)
	assert.Equal(t, "test-base", context.Git.Base)
}

func TestLoadTemplateStringsErrorPublishStrategy(t *testing.T) {
	context := &ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{},
			Commit:  &v1.CommitConfig{},
		},
		PublishStrategy: strategies.PublishStrategy("invalid"),
	}

	err := loadTemplateStrings(context)
	assert.EqualError(t, err, "unsupported publish strategy: invalid")
}

func TestLoadReleaseConstraints(t *testing.T) {
	context := &ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{
				Matches: []string{".*"},
				Ignores: []string{"master"},
			},
		},
	}

	err := loadReleaseConstraints(context)
	assert.NoError(t, err)

	assert.Len(t, context.Release.Matches, 1)
	assert.Equal(t, ".*", context.Release.Matches[0].String())
	assert.Len(t, context.Release.Ignores, 1)
	assert.Equal(t, "master", context.Release.Ignores[0].String())
}

func TestLoadReleaseConstraintsMatchError(t *testing.T) {
	context := &ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{
				Matches: []string{"*"},
			},
		},
	}

	err := loadReleaseConstraints(context)
	assert.EqualError(t, err, "error parsing regexp: missing argument to repetition operator: `*`")

	assert.Len(t, context.Release.Matches, 0)
	assert.Len(t, context.Release.Ignores, 0)
}

func TestLoadReleaseConstraintsMatchErrorDryRun(t *testing.T) {
	context := &ctx.Context{
		DryRun: true,
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{
				Matches: []string{"*"},
				Ignores: []string{"master"},
			},
		},
	}

	err := loadReleaseConstraints(context)
	assert.NoError(t, err)

	assert.Len(t, context.Release.Matches, 1)
	assert.Equal(t, "dry-run", context.Release.Matches[0].String())
	assert.Len(t, context.Release.Ignores, 1)
	assert.Equal(t, "master", context.Release.Ignores[0].String())
	assert.EqualError(t, context.Errors(), "\nErrors:\n • error parsing regexp: missing argument to repetition operator: `*`\n\n")
}

func TestLoadReleaseConstraintsIgnoreError(t *testing.T) {
	context := &ctx.Context{
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{
				Ignores: []string{"*"},
			},
		},
	}

	err := loadReleaseConstraints(context)
	assert.EqualError(t, err, "error parsing regexp: missing argument to repetition operator: `*`")

	assert.Len(t, context.Release.Matches, 0)
	assert.Len(t, context.Release.Ignores, 0)
}

func TestLoadReleaseConstraintsIgnoreErrorDryRun(t *testing.T) {
	context := &ctx.Context{
		DryRun: true,
		Config: &v1.Config{
			Release: &v1.ReleaseConfig{
				Matches: []string{"master"},
				Ignores: []string{"*"},
			},
		},
	}

	err := loadReleaseConstraints(context)
	assert.NoError(t, err)

	assert.Len(t, context.Release.Matches, 1)
	assert.Equal(t, "master", context.Release.Matches[0].String())
	assert.Len(t, context.Release.Ignores, 1)
	assert.Equal(t, "dry-run", context.Release.Ignores[0].String())
	assert.EqualError(t, context.Errors(), "\nErrors:\n • error parsing regexp: missing argument to repetition operator: `*`\n\n")
}
