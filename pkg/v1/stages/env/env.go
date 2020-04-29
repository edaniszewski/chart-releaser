package env

import (
	"errors"
	"os"

	"github.com/apex/log"
	"github.com/edaniszewski/chart-releaser/pkg/env"
	context "github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
)

// Errors for the env stage.
var (
	ErrGithubTokenNotSet   = errors.New("GITHUB_TOKEN environment variable not set")
	ErrRepoTypeNotSet      = errors.New("repository type not set prior to running 'env' stage")
	ErrUnsupportedRepoType = errors.New("unsupported repository type")
)

// Stage for the "env" step of the update pipeline.
type Stage struct{}

// Name of the stage.
func (Stage) Name() string {
	return "env"
}

// String describes what the stage does.
func (Stage) String() string {
	return "loading environment variables"
}

// Run the operations defined for the stage.
func (Stage) Run(ctx *context.Context) error {
	if ctx.Repository.Type == "" {
		return ErrRepoTypeNotSet
	}

	switch ctx.Repository.Type {
	case context.RepoGithub:
		val, found := os.LookupEnv(env.GithubToken)
		if !found {
			if err := ctx.CheckDryRun(ErrGithubTokenNotSet); err != nil {
				return err
			}
			val = ""
			log.Warn("github token not detected - using no token for dry-run")
		}
		ctx.Token = val

	default:
		log.WithFields(log.Fields{
			"type": ctx.Repository.Type,
		}).Error("unsupported repository type specified")
		return ErrUnsupportedRepoType
	}
	return nil
}
