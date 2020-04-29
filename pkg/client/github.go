package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

// githubClient implements the Client interface for updating charts
// on GitHub.
type githubClient struct {
	client *github.Client
}

// NewGitHubClient creates a new GitHub client.
func NewGitHubClient(ctx context.Context, token string) (Client, error) {
	if token == "" {
		return nil, fmt.Errorf("no token provided to github client")
	}

	tok := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	client := github.NewClient(oauth2.NewClient(ctx, tok))

	return githubClient{
		client: client,
	}, nil
}

func verifyOptions(opts *Options) error {
	if !strings.HasPrefix(opts.Ref, "refs/heads/") {
		opts.Ref = fmt.Sprintf("refs/heads/%s", opts.Ref)
	}

	if !strings.HasPrefix(opts.Base, "refs/heads/") {
		opts.Base = fmt.Sprintf("refs/heads/%s", opts.Base)
	}
	return nil
}

// GetFile gets the data for the specified file from the chart repository.
func (c githubClient) GetFile(ctx context.Context, opts *Options, path string) (string, error) {
	if err := verifyOptions(opts); err != nil {
		return "", err
	}

	file, _, _, err := c.client.Repositories.GetContents(
		ctx,
		opts.RepoOwner,
		opts.RepoName,
		path,
		&github.RepositoryContentGetOptions{},
	)
	if err != nil {
		return "", err
	}

	contents, err := file.GetContent()
	if err != nil {
		return "", err
	}

	return contents, nil
}

// UpdateFile updates the content of the specified file within the configured
// chart repository.
func (c githubClient) UpdateFile(ctx context.Context, opts *Options, path string, msg string, contents []byte) error {
	if err := verifyOptions(opts); err != nil {
		return err
	}

	options := &github.RepositoryContentFileOptions{
		Committer: &github.CommitAuthor{
			Name:  github.String(opts.AuthorName),
			Email: github.String(opts.AuthorEmail),
		},
		Content: contents,
		Message: github.String(msg),
		Branch:  github.String(opts.Ref),
	}

	// First, check that the file exists in the specified repo.
	file, _, resp, err := c.client.Repositories.GetContents(
		ctx,
		opts.RepoOwner,
		opts.RepoName,
		path,
		&github.RepositoryContentGetOptions{},
	)
	if err != nil {
		// If we get a 404, the file does not exist. Return an error. It is up to
		// the caller to decide whether or not to create the file or to error out.
		if resp.StatusCode == 404 {
			log.WithFields(log.Fields{
				"error": err,
				"file":  path,
				"ref":   opts.Ref,
			}).Error("github client: unable to update file (not found)")
			return ErrFileNotFound
		}
		return err
	}

	// The file exists -- update it.
	options.SHA = file.SHA
	_, _, err = c.client.Repositories.UpdateFile(
		ctx,
		opts.RepoOwner,
		opts.RepoName,
		path,
		options,
	)
	return err
}

// CreateRef creates a new ref to stage the commits produced by chart-releaser.
// If the ref already exists, an error is returned.
func (c githubClient) CreateRef(ctx context.Context, opts *Options) error {
	if err := verifyOptions(opts); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"repo": fmt.Sprintf("%s/%s", opts.RepoOwner, opts.RepoName),
		"ref":  opts.Base,
	}).Debug("github client: getting reference")
	ref, _, err := c.client.Git.GetRef(
		ctx,
		opts.RepoOwner,
		opts.RepoName,
		opts.Base,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"base":  opts.Base,
			"ref":   opts.Ref,
		}).Error("github client: unable to create ref - configured base ref does not exist")
		return err
	}

	log.WithFields(log.Fields{
		"repo": fmt.Sprintf("%s/%s", opts.RepoOwner, opts.RepoName),
		"ref":  opts.Ref,
	}).Debug("github client: creating reference")
	_, _, err = c.client.Git.CreateRef(
		ctx,
		opts.RepoOwner,
		opts.RepoName,
		&github.Reference{
			Ref: &opts.Ref,
			Object: &github.GitObject{
				SHA: ref.Object.SHA,
			},
		},
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ref":   opts.Ref,
			"base":  opts.Base,
		}).Error("github client: failed to create new ref")
		return err
	}
	return nil
}

// CreatePullRequest creates a new pull request for the changes produced by chart-releaser.
func (c githubClient) CreatePullRequest(ctx context.Context, opts *Options, title, body string) error {
	if err := verifyOptions(opts); err != nil {
		return err
	}

	if opts.Base == opts.Ref {
		// todo: logging
		return fmt.Errorf("cannot create pull request, ref and base are the same")
	}

	log.WithFields(log.Fields{
		"repo":  fmt.Sprintf("%s/%s", opts.RepoOwner, opts.RepoName),
		"ref":   opts.Ref,
		"base":  opts.Base,
		"title": title,
	}).Debug("github client: creating pull request")
	pr, _, err := c.client.PullRequests.Create(
		ctx,
		opts.RepoOwner,
		opts.RepoName,
		&github.NewPullRequest{
			Title: &title,
			Body:  &body,
			Base:  &opts.Base,
			Head:  &opts.Ref,
		},
	)
	if err != nil {
		// todo: logging
		return err
	}

	log.Infof("created pull request %v %v (%v <- %v)", pr.GetNumber(), pr.GetURL(), pr.GetBase().GetRef(), pr.GetHead().GetRef())
	return nil
}
