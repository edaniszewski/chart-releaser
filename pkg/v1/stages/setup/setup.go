package setup

import (
	"errors"

	"github.com/apex/log"
	"github.com/edaniszewski/chart-releaser/pkg/utils"
	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
)

// Errors for the setup stage.
var (
	ErrGitNotFound = errors.New("git not found on PATH")
	ErrNotInRepo   = errors.New("current directory is not a git repository")
	ErrDirtyGit    = errors.New("git is in a dirty state")
)

// Stage for the "setup" step of the update pipeline.
type Stage struct{}

// Name of the stage.
func (Stage) Name() string {
	return "setup"
}

// String describes what the stage does.
func (Stage) String() string {
	return "performing pre-flight setup and checks"
}

// Run the operations defined for the stage.
func (Stage) Run(ctx *ctx.Context) error {
	log.Debug("checking if git exists on PATH")
	if !utils.BinExists("git") {
		return ErrGitNotFound
	}

	log.Debug("checking if directory is a git repo")
	if !utils.InRepo() {
		if err := ctx.CheckDryRun(ErrNotInRepo); err != nil {
			return err
		}
	}

	log.Debug("checking if git is in a clean state")
	if isDirty, out := utils.IsDirty(); isDirty {
		if ctx.AllowDirty {
			log.Info("allowing git to be in a dirty state")
		} else {
			if err := ctx.CheckDryRun(ErrDirtyGit); err != nil {
				log.Errorf("dirty git state detected\n" + out)
				return err
			}
		}
	}

	return nil
}
