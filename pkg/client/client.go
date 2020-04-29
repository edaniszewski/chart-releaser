package client

import (
	"context"
	"errors"
)

// Errors relating to client operations.
var (
	ErrFileNotFound = errors.New("file not found in remote repo")
)

// The Client interface defines a way to interact with a source repository
// to be able to perform operations on Helm Charts and other project files.
type Client interface {
	GetFile(ctx context.Context, opts *Options, path string) (contents string, err error)
	UpdateFile(ctx context.Context, opts *Options, path string, msg string, contents []byte) error
	CreateRef(ctx context.Context, opts *Options) error
	CreatePullRequest(ctx context.Context, opts *Options, title, body string) error
}

// Options are the configuration options and state required to create a new
// client.
type Options struct {
	Ref         string
	Base        string
	RepoName    string
	RepoOwner   string
	AuthorName  string
	AuthorEmail string
}
