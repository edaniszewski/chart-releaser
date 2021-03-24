package ctx

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/edaniszewski/chart-releaser/pkg/errs"
	v1 "github.com/edaniszewski/chart-releaser/pkg/v1/cfg"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestListRepoTypes(t *testing.T) {
	types := ListRepoTypes()
	assert.Len(t, types, 1)
}

func TestRepoTypeFromString(t *testing.T) {
	rt, err := RepoTypeFromString("github")
	assert.NoError(t, err)
	assert.Equal(t, RepoGithub, rt)

	rt, err = RepoTypeFromString("github.com")
	assert.NoError(t, err)
	assert.Equal(t, RepoGithub, rt)
}

func TestRepoTypeFromString_Error(t *testing.T) {
	_, err := RepoTypeFromString("invalid")
	assert.EqualError(t, err, "unsupported repository type: invalid (supported: [github])")
}

func TestFile_HasChanges_Empty(t *testing.T) {
	file := File{
		PreviousContents: []byte{},
		NewContents:      []byte{},
	}
	assert.False(t, file.HasChanges())
}

func TestFile_HasChanges_NoChanges(t *testing.T) {
	file := File{
		PreviousContents: []byte{0x01, 0x02, 0x03},
		NewContents:      []byte{0x01, 0x02, 0x03},
	}
	assert.False(t, file.HasChanges())
}

func TestFile_HasChanges(t *testing.T) {
	file := File{
		PreviousContents: []byte{0x01, 0x02, 0x03},
		NewContents:      []byte{0x02, 0x03, 0x04},
	}
	assert.True(t, file.HasChanges())
}

func TestNew(t *testing.T) {
	c := v1.Config{}

	ctx := New(&c)
	assert.NotNil(t, ctx.Context)
	assert.Equal(t, &c, ctx.Config)
}

func TestNewWithTimeout(t *testing.T) {
	c := v1.Config{}

	ctx, _ := NewWithTimeout(&c, 5*time.Second)
	assert.NotNil(t, ctx.Context)
	assert.Equal(t, &c, ctx.Config)
}

func TestWrap(t *testing.T) {
	c := v1.Config{}
	ctxt := context.Background()

	ctx := Wrap(ctxt, &c)
	assert.Equal(t, ctxt, ctx.Context)
	assert.Equal(t, &c, ctx.Config)
}

func TestContext_PrintErrors_NoErrors(t *testing.T) {
	buf := bytes.Buffer{}
	ctx := Context{
		Out:    &buf,
		errors: errs.Collector{},
	}

	ctx.PrintErrors()
	assert.Equal(t, "\x1b[0;32mdry-run completed without errors\n\x1b[0m", buf.String())
}

func TestContext_PrintErrors_HasErrors(t *testing.T) {
	buf := bytes.Buffer{}
	ctx := Context{
		Out:    &buf,
		errors: errs.Collector{},
	}
	ctx.errors.Add(errors.New("test error"))

	ctx.PrintErrors()
	assert.Equal(t, "\x1b[0;31m\ndry-run completed with errors\x1b[0m\nErrors:\n • test error\n\n", buf.String())
}

func TestContext_CheckDryRun(t *testing.T) {
	ctx := Context{
		DryRun: true,
		errors: errs.Collector{},
	}

	err := ctx.CheckDryRun(errors.New("test error"))
	assert.NoError(t, err)
	assert.Equal(t, 1, ctx.errors.Count())
}

func TestContext_CheckDryRun_NotDryRun(t *testing.T) {
	ctx := Context{
		DryRun: false,
		errors: errs.Collector{},
	}

	err := ctx.CheckDryRun(errors.New("test error"))
	assert.EqualError(t, err, "test error")
	assert.Equal(t, 0, ctx.errors.Count())
}

func TestContext_HasErrors(t *testing.T) {
	ctx := Context{
		errors: errs.Collector{},
	}
	ctx.errors.Add(errors.New("test error"))

	assert.EqualError(t, ctx.Errors(), "\nErrors:\n • test error\n\n")
}

func TestContext_NoErrors(t *testing.T) {
	ctx := Context{
		errors: errs.Collector{},
	}

	assert.Nil(t, ctx.Errors())
}
