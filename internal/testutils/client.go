package testutils

import (
	"context"

	"github.com/edaniszewski/chart-releaser/pkg/client"
)

// FakeClient implements the Client interface. It is used for testing.
type FakeClient struct {
	FileData string

	GetFileError           []error
	UpdateFileError        []error
	CreateRefError         []error
	CreatePullRequestError []error

	getIdx       int
	updateIdx    int
	createRefIdx int
	createPRIdx  int
}

func (c *FakeClient) GetFile(ctx context.Context, opts *client.Options, path string) (string, error) {
	if len(c.GetFileError) == 0 {
		return c.FileData, nil
	}
	data := c.GetFileError[c.getIdx]
	c.getIdx++
	return c.FileData, data
}

func (c *FakeClient) UpdateFile(ctx context.Context, opts *client.Options, path string, msg string, contents []byte) error {
	if len(c.UpdateFileError) == 0 {
		return nil
	}
	data := c.UpdateFileError[c.updateIdx]
	c.updateIdx++
	return data
}

func (c *FakeClient) CreateRef(ctx context.Context, opts *client.Options) error {
	if len(c.CreateRefError) == 0 {
		return nil
	}
	data := c.CreateRefError[c.createRefIdx]
	c.createRefIdx++
	return data
}

func (c *FakeClient) CreatePullRequest(ctx context.Context, opts *client.Options, title, body string) error {
	if len(c.CreatePullRequestError) == 0 {
		return nil
	}
	data := c.CreatePullRequestError[c.createPRIdx]
	c.createPRIdx++
	return data
}
