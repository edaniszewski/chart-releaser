package ctx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/davecgh/go-spew/spew"
	"github.com/edaniszewski/chart-releaser/pkg/client"
	"github.com/edaniszewski/chart-releaser/pkg/errs"
	version "github.com/edaniszewski/chart-releaser/pkg/semver"
	"github.com/edaniszewski/chart-releaser/pkg/strategies"
	v1 "github.com/edaniszewski/chart-releaser/pkg/v1/cfg"
	"github.com/fatih/color"
	"github.com/mgutz/ansi"
)

// RepoType defines the type of repository container the Helm Chart.
type RepoType string

// RepoTypes supported by chart-releaser
const (
	RepoGithub RepoType = "github"
)

// ListRepoTypes returns all of the RepoTypes supported by chart-releaser.
func ListRepoTypes() []RepoType {
	return []RepoType{
		RepoGithub,
	}
}

// RepoTypeFromString gets a RepoType corresponding to a string.
func RepoTypeFromString(s string) (RepoType, error) {
	switch strings.ToLower(s) {
	case "github", "github.com":
		return RepoGithub, nil
	default:
		return "", fmt.Errorf("unsupported repository type: %s (supported: %v)", s, ListRepoTypes())
	}
}

// App version information.
type App struct {
	NewVersion      version.Semver
	PreviousVersion version.Semver
}

// Author information for the committer.
type Author struct {
	Name  string
	Email string
}

// Chart version, metadata, and contents information.
type Chart struct {
	Name            string
	SubPath         string
	File            File
	NewVersion      version.Semver
	PreviousVersion version.Semver
}

// File holds basic data about a file, its previous contents,
// and its updated contents.
type File struct {
	Path             string
	PreviousContents []byte
	NewContents      []byte
}

// HasChanges determines whether the previous contents of the file differ from the
// new contents of the file.
func (f *File) HasChanges() bool {
	return !bytes.Equal(f.PreviousContents, f.NewContents)
}

// Git information used for publishing chart updates.
type Git struct {
	Tag  string
	Ref  string
	Base string
}

// Repository metadata.
type Repository struct {
	Type  RepoType
	Owner string
	Name  string
}

// Release metadata used for generating the release messages (commits, PRs).
type Release struct {
	PRTitle         string
	PRBody          string
	UpdateCommitMsg string
	Matches         []*regexp.Regexp
	Ignores         []*regexp.Regexp
}

// Context holds information that is used by chart-releaser throughout
// the update process. Its values get populated and updated as various
// stages operate on it. It holds all release state.
type Context struct {
	context.Context
	Config *v1.Config
	Out    io.Writer

	Token           string
	PublishStrategy strategies.PublishStrategy
	UpdateStrategy  strategies.UpdateStrategy

	App        App
	Author     Author
	Chart      Chart
	Client     client.Client
	Files      []File
	Git        Git
	Repository Repository
	Release    Release

	AllowDirty bool
	DryRun     bool
	ShowDiff   bool

	errors errs.Collector
}

// Dump the Context to console.
func (ctx *Context) Dump() {
	if ctx.Out == nil {
		log.Error("unable to dump context: context output writer is nil")
		return
	}
	_, _ = fmt.Fprintln(ctx.Out, color.New(color.Bold).Sprint("\n=== Context ==="))
	_, _ = fmt.Fprintln(ctx.Out, fmt.Sprintf("AllowDirty:\t\t%v", ctx.AllowDirty))
	_, _ = fmt.Fprintln(ctx.Out, fmt.Sprintf("DryRun:\t\t\t%v", ctx.DryRun))
	_, _ = fmt.Fprintln(ctx.Out, fmt.Sprintf("ShowDiff:\t\t%v", ctx.ShowDiff))
	_, _ = fmt.Fprintln(ctx.Out, fmt.Sprintf("PublishStrategy:\t%s", ctx.PublishStrategy))
	_, _ = fmt.Fprintln(ctx.Out, fmt.Sprintf("UpdateStrategy:\t\t%s", ctx.UpdateStrategy))
	_, _ = fmt.Fprintln(ctx.Out, fmt.Sprintf("Token:\t\t\t%s****", ctx.Token[0:4]))
	_, _ = fmt.Fprintln(ctx.Out, color.New(color.Bold).Sprint("Config"))
	_, _ = fmt.Fprintln(ctx.Out, spew.Sdump(ctx.Config))
	_, _ = fmt.Fprintln(ctx.Out, color.New(color.Bold).Sprint("App"))
	_, _ = fmt.Fprintln(ctx.Out, spew.Sdump(ctx.App))
	_, _ = fmt.Fprintln(ctx.Out, color.New(color.Bold).Sprint("Author"))
	_, _ = fmt.Fprintln(ctx.Out, spew.Sdump(ctx.Author))
	_, _ = fmt.Fprintln(ctx.Out, color.New(color.Bold).Sprint("Chart"))
	_, _ = fmt.Fprintln(ctx.Out, spew.Sdump(ctx.Chart))
	_, _ = fmt.Fprintln(ctx.Out, color.New(color.Bold).Sprint("Files"))
	_, _ = fmt.Fprintln(ctx.Out, spew.Sdump(ctx.Files))
	_, _ = fmt.Fprintln(ctx.Out, color.New(color.Bold).Sprint("Git"))
	_, _ = fmt.Fprintln(ctx.Out, spew.Sdump(ctx.Git))
	_, _ = fmt.Fprintln(ctx.Out, color.New(color.Bold).Sprint("Repository"))
	_, _ = fmt.Fprintln(ctx.Out, spew.Sdump(ctx.Repository))
	_, _ = fmt.Fprintln(ctx.Out, color.New(color.Bold).Sprint("Release"))
	_, _ = fmt.Fprintln(ctx.Out, spew.Sdump(ctx.Release))
}

// PrintErrors prints any errors encountered and collected by the Context's error
// collector.
func (ctx *Context) PrintErrors() {
	if ctx.errors.HasErrors() {
		_, _ = fmt.Fprint(ctx.Out, ansi.Color("\ndry-run completed with errors", "red"))
		_, _ = fmt.Fprint(ctx.Out, ctx.errors.Error())
	} else {
		_, _ = fmt.Fprint(ctx.Out, ansi.Color("dry-run completed without errors\n", "green"))
	}
}

// Errors returns the errors collected by the Context.
func (ctx *Context) Errors() error {
	if ctx.errors.HasErrors() {
		return &ctx.errors
	}
	return nil
}

// CheckDryRun checks whether the Context is configured for a dry-run, and if so,
// to only log the provided error. If it is not in a dry-run state, it will return
// the same error it was given for the caller to propagate appropriately.
func (ctx *Context) CheckDryRun(err error) error {
	if ctx.DryRun {
		log.WithError(err).Warn("dry-run: ignoring error")
		ctx.errors.Add(err)
		return nil
	}
	return err
}

// New creates a new v1 Context for the given v1 Config.
func New(config *v1.Config) *Context {
	return Wrap(context.Background(), config)
}

// NewWithTimeout creates a new v1 Context for the given v1 Config, with
// an timeout to cancel the context.
func NewWithTimeout(config *v1.Config, timeout time.Duration) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return Wrap(ctx, config), cancel
}

// Wrap an existing context.Context into a v1 Context.
func Wrap(ctx context.Context, config *v1.Config) *Context {
	return &Context{
		Context: ctx,
		Config:  config,
		Out:     os.Stdout,
		errors:  *errs.NewCollector(),
	}
}
