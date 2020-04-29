package client

import (
	"errors"

	"github.com/apex/log"
	"github.com/edaniszewski/chart-releaser/pkg/client"
	context "github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
)

// Errors for the client stage.
var (
	ErrRepoTypeNotSet      = errors.New("repository type not set prior to running 'client' stage")
	ErrUnsupportedRepoType = errors.New("unsupported repository type")
)

// Stage for the "client" step of the update pipeline.
type Stage struct{}

// Name of the stage.
func (Stage) Name() string {
	return "client"
}

// String describes what the stage does.
func (Stage) String() string {
	return "creating repository client"
}

// Run the operations defined for the stage.
func (Stage) Run(ctx *context.Context) error {
	if ctx.Repository.Type == "" {
		return ErrRepoTypeNotSet
	}

	switch ctx.Repository.Type {
	case context.RepoGithub:
		c, err := client.NewGitHubClient(ctx.Context, ctx.Token)
		if err != nil {
			return err
		}
		ctx.Client = c

	default:
		log.WithFields(log.Fields{
			"type": ctx.Repository.Type,
		}).Error("unsupported repository type specified")
		return ErrUnsupportedRepoType
	}
	return nil
}
