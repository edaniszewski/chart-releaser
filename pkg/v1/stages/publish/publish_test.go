package publish

import (
	"regexp"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/edaniszewski/chart-releaser/internal/testutils"
	"github.com/edaniszewski/chart-releaser/pkg/strategies"
	"github.com/edaniszewski/chart-releaser/pkg/templates"
	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestStage_Name(t *testing.T) {
	assert.Equal(t, "publish", Stage{}.Name())
}

func TestStage_String(t *testing.T) {
	assert.Equal(t, "publishing changes", Stage{}.String())
}

func TestStage_Run_StrategyCommit(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{},
		Git: ctx.Git{
			Ref:  "ref-branch",
			Base: "base-branch",
		},
		Release: ctx.Release{
			UpdateCommitMsg: "update-msg",
			PRTitle:         "pr-title",
			PRBody:          "pr-body",
		},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
		PublishStrategy: strategies.PublishCommit,
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "base-branch", context.Git.Base)
	assert.Equal(t, "update-msg", context.Release.UpdateCommitMsg)
	assert.Equal(t, "pr-title", context.Release.PRTitle)
	assert.Equal(t, "pr-body", context.Release.PRBody)
}

func TestStage_Run_StrategyPR(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{},
		Git: ctx.Git{
			Ref:  "ref-branch",
			Base: "base-branch",
		},
		Release: ctx.Release{
			UpdateCommitMsg: "update-msg",
			PRTitle:         "pr-title",
			PRBody:          "pr-body",
		},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
		PublishStrategy: strategies.PublishPullRequest,
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "base-branch", context.Git.Base)
	assert.Equal(t, "update-msg", context.Release.UpdateCommitMsg)
	assert.Equal(t, "pr-title", context.Release.PRTitle)
	assert.Equal(t, "pr-body", context.Release.PRBody)
}

func TestStage_Run_StrategyPRDefaultTemplates(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{},
		Git: ctx.Git{
			Ref:  templates.DefaultBranchName,
			Base: "master",
		},
		Release: ctx.Release{
			UpdateCommitMsg: templates.DefaultUpdateCommitMessage,
			PRTitle:         templates.DefaultPullRequestTitle,
			PRBody:          templates.DefaultPullRequestBody,
		},
		Chart: ctx.Chart{
			Name:            "test-chart",
			NewVersion:      testutils.NewSemver(t, "1.2.3"),
			PreviousVersion: testutils.NewSemver(t, "1.2.2"),
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "extra/file/1.txt",
				PreviousContents: []byte("old"),
				NewContents:      []byte("new"),
			},
			{
				Path:             "extra/file/2.txt",
				PreviousContents: []byte("old"),
				NewContents:      []byte("new"),
			},
		},
		App: ctx.App{
			NewVersion:      testutils.NewSemver(t, "1.0.0"),
			PreviousVersion: testutils.NewSemver(t, "0.1.0"),
		},
		PublishStrategy: strategies.PublishPullRequest,
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)

	assert.Equal(t, "chartreleaser/test-chart/1.2.3", context.Git.Ref)
	assert.Equal(t, "master", context.Git.Base)
	assert.Equal(t, "[test-chart] bump chart to 1.2.3 for new application release (1.0.0)", context.Release.UpdateCommitMsg)
	assert.Equal(t, "Bump test-chart Chart from 1.2.2 to 1.2.3", context.Release.PRTitle)
	assert.Equal(t, heredoc.Doc(`
		Bumps the test-chart Helm Chart from 1.2.2 to 1.2.3.

		The following files have also been updated:
		- extra/file/1.txt
		- extra/file/2.txt
		
		---
		*This PR was generated with [chart-releaser](https://github.com/edaniszewski/chart-releaser)*
`,
	), context.Release.PRBody)
}

func TestStage_Run_DryRun(t *testing.T) {
	context := ctx.Context{
		DryRun: true,
		Git: ctx.Git{
			Ref:  "ref-branch",
			Base: "base-branch",
		},
		Release: ctx.Release{
			UpdateCommitMsg: "update-msg",
			PRTitle:         "pr-title",
			PRBody:          "pr-body",
		},
		PublishStrategy: strategies.PublishCommit,
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "base-branch", context.Git.Base)
	assert.Equal(t, "update-msg", context.Release.UpdateCommitMsg)
	assert.Equal(t, "pr-title", context.Release.PRTitle)
	assert.Equal(t, "pr-body", context.Release.PRBody)
}

func TestStage_RunPublishStrategyError(t *testing.T) {
	context := ctx.Context{
		Git: ctx.Git{
			Ref:  "ref-branch",
			Base: "base-branch",
		},
		Release: ctx.Release{
			UpdateCommitMsg: "update-msg",
			PRTitle:         "pr-title",
			PRBody:          "pr-body",
		},
		PublishStrategy: strategies.PublishStrategy("invalid"),
	}

	err := Stage{}.Run(&context)
	assert.Equal(t, ErrUnsupportedPublishStrategy, err)

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "base-branch", context.Git.Base)
	assert.Equal(t, "update-msg", context.Release.UpdateCommitMsg)
	assert.Equal(t, "pr-title", context.Release.PRTitle)
	assert.Equal(t, "pr-body", context.Release.PRBody)
}

func TestStage_RunNoTagMatch(t *testing.T) {
	context := ctx.Context{
		Git: ctx.Git{
			Ref:  "ref-branch",
			Tag:  "dev",
			Base: "base-branch",
		},
		Release: ctx.Release{
			UpdateCommitMsg: "update-msg",
			PRTitle:         "pr-title",
			PRBody:          "pr-body",
			Matches: []*regexp.Regexp{
				regexp.MustCompile(`(\.[0-9])+`),
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "base-branch", context.Git.Base)
	assert.Equal(t, "update-msg", context.Release.UpdateCommitMsg)
	assert.Equal(t, "pr-title", context.Release.PRTitle)
	assert.Equal(t, "pr-body", context.Release.PRBody)
}

func TestStage_RunNoTagMatchDryRun(t *testing.T) {
	context := ctx.Context{
		DryRun: true,
		Client: &testutils.FakeClient{},
		Git: ctx.Git{
			Ref:  "ref-branch",
			Tag:  "dev",
			Base: "base-branch",
		},
		Release: ctx.Release{
			UpdateCommitMsg: "update-msg",
			PRTitle:         "pr-title",
			PRBody:          "pr-body",
			Matches: []*regexp.Regexp{
				regexp.MustCompile(`(\.[0-9])+`),
			},
		},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "base-branch", context.Git.Base)
	assert.Equal(t, "update-msg", context.Release.UpdateCommitMsg)
	assert.Equal(t, "pr-title", context.Release.PRTitle)
	assert.Equal(t, "pr-body", context.Release.PRBody)
}

func TestStage_RunTagIgnoreMatch(t *testing.T) {
	context := ctx.Context{
		Git: ctx.Git{
			Ref:  "ref-branch",
			Tag:  "dev",
			Base: "base-branch",
		},
		Release: ctx.Release{
			UpdateCommitMsg: "update-msg",
			PRTitle:         "pr-title",
			PRBody:          "pr-body",
			Ignores: []*regexp.Regexp{
				regexp.MustCompile(`dev(.)*`),
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "base-branch", context.Git.Base)
	assert.Equal(t, "update-msg", context.Release.UpdateCommitMsg)
	assert.Equal(t, "pr-title", context.Release.PRTitle)
	assert.Equal(t, "pr-body", context.Release.PRBody)
}

func TestStage_RunTagIgnoreMatchDryRun(t *testing.T) {
	context := ctx.Context{
		DryRun: true,
		Client: &testutils.FakeClient{},
		Git: ctx.Git{
			Ref:  "ref-branch",
			Tag:  "dev",
			Base: "base-branch",
		},
		Release: ctx.Release{
			UpdateCommitMsg: "update-msg",
			PRTitle:         "pr-title",
			PRBody:          "pr-body",
			Ignores: []*regexp.Regexp{
				regexp.MustCompile(`dev(.)*`),
			},
		},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
	}

	err := Stage{}.Run(&context)
	assert.NoError(t, err)

	assert.Equal(t, "ref-branch", context.Git.Ref)
	assert.Equal(t, "base-branch", context.Git.Base)
	assert.Equal(t, "update-msg", context.Release.UpdateCommitMsg)
	assert.Equal(t, "pr-title", context.Release.PRTitle)
	assert.Equal(t, "pr-body", context.Release.PRBody)
}

func Test_PublishCommit(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("hello"),
				NewContents:      []byte("hi"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "extra1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
			{
				Path:             "extra2",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("foo"),
			},
		},
	}

	err := publishCommit(&context)
	assert.NoError(t, err)
}

func Test_PublishCommitCreateRefError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{
			CreateRefError: []error{
				errors.New("test error"),
			},
		},
		Git: ctx.Git{
			Ref: "testbranch",
		},
	}

	err := publishCommit(&context)
	assert.EqualError(t, err, "test error")
}

func Test_PublishCommitUpdateFileError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{
			UpdateFileError: []error{
				errors.New("test error"),
			},
		},
		Git: ctx.Git{
			Ref: "testbranch",
		},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("hello"),
				NewContents:      []byte("hi"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "extra1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
	}

	err := publishCommit(&context)
	assert.EqualError(t, err, "test error")
}

func Test_PublishCommitNoChartChanges(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{},
		Git: ctx.Git{
			Ref: "testbranch",
		},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("hello"),
				NewContents:      []byte("hello"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "extra1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
	}

	err := publishCommit(&context)
	assert.Equal(t, ErrNoChartChanges, err)
}

func Test_PublishCommitExtrasUpdateFileError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{
			UpdateFileError: []error{
				nil,
				errors.New("test error"),
			},
		},
		Git: ctx.Git{
			Ref: "testbranch",
		},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("hello"),
				NewContents:      []byte("hi"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "extra1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
	}

	err := publishCommit(&context)
	assert.EqualError(t, err, "test error")
}

func Test_PublishPullRequest(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("hello"),
				NewContents:      []byte("hi"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "extra1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
	}

	err := publishPullRequest(&context)
	assert.NoError(t, err)

}

func Test_PublishPullRequest_PublishCommitError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{
			CreateRefError: []error{
				errors.New("test error"),
			},
		},
		Git: ctx.Git{
			Ref: "testbranch",
		},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("hello"),
				NewContents:      []byte("hi"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "extra1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
	}

	err := publishPullRequest(&context)
	assert.EqualError(t, err, "test error")
}

func Test_PublishPullRequestParseTitleError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("hello"),
				NewContents:      []byte("hi"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "extra1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
		Release: ctx.Release{
			PRTitle: "{{",
		},
	}

	err := publishPullRequest(&context)
	assert.EqualError(t, err, "template: :1: unexpected unclosed action in command")
}

func Test_PublishPullRequestExecuteTitleError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("hello"),
				NewContents:      []byte("hi"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "extra1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
		Release: ctx.Release{
			PRTitle: "{{ .Not.Exists }}",
		},
	}

	err := publishPullRequest(&context)
	assert.EqualError(t, err, "template: :1:7: executing \"\" at <.Not.Exists>: can't evaluate field Not in type *ctx.Context")
}

func Test_PublishPullRequestParseBodyError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("hello"),
				NewContents:      []byte("hi"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "extra1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
		Release: ctx.Release{
			PRBody: "{{",
		},
	}

	err := publishPullRequest(&context)
	assert.EqualError(t, err, "template: :1: unexpected unclosed action in command")
}

func Test_PublishPullRequestExecuteBodyError(t *testing.T) {
	context := ctx.Context{
		Client: &testutils.FakeClient{},
		Chart: ctx.Chart{
			File: ctx.File{
				Path:             "path1",
				PreviousContents: []byte("hello"),
				NewContents:      []byte("hi"),
			},
		},
		Files: []ctx.File{
			{
				Path:             "extra1",
				PreviousContents: []byte("foo"),
				NewContents:      []byte("bar"),
			},
		},
		Release: ctx.Release{
			PRBody: "{{ .Not.Exists }}",
		},
	}

	err := publishPullRequest(&context)
	assert.EqualError(t, err, "template: :1:7: executing \"\" at <.Not.Exists>: can't evaluate field Not in type *ctx.Context")
}
