package v1

import (
	"testing"

	"github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
	"github.com/stretchr/testify/assert"
)

func TestConfigVersion(t *testing.T) {
	assert.Equal(t, "v1", ConfigVersion())
}

func TestNewUpdater(t *testing.T) {
	u := NewUpdater([]byte{0x00, 0x01})
	assert.NotNil(t, u)
	assert.Equal(t, []byte{0x00, 0x01}, u.data)
}

func TestUpdater_Run(t *testing.T) {
	//  todo
}

func TestUpdater_Run_LoadError(t *testing.T) {
	u := NewUpdater([]byte{0x00, 0x01})

	err := u.Run(UpdateOptions{})
	assert.EqualError(t, err, "error converting YAML to JSON: yaml: control characters are not allowed")
}

func TestUpdateOptions_AugmentCtx(t *testing.T) {
	context := ctx.Context{}
	assert.False(t, context.AllowDirty)
	assert.False(t, context.DryRun)
	assert.False(t, context.ShowDiff)

	opts := UpdateOptions{
		AllowDirty: true,
		ShowDiff:   true,
		DryRun:     true,
	}

	opts.AugmentCtx(&context)

	assert.True(t, context.AllowDirty)
	assert.True(t, context.DryRun)
	assert.True(t, context.ShowDiff)
}

func TestNewFormatter(t *testing.T) {
	f := NewFormatter([]byte{0x00, 0x01})
	assert.NotNil(t, f)
	assert.Equal(t, []byte{0x00, 0x01}, f.data)
}

func TestFormatter_Run(t *testing.T) {
	//  todo
}

func TestFormatter_Run_LoadError(t *testing.T) {
	f := NewFormatter([]byte{0x00, 0x01})

	err := f.Run(FormatterOptions{})
	assert.EqualError(t, err, "error converting YAML to JSON: yaml: control characters are not allowed")
}

func TestNewChecker(t *testing.T) {
	c := NewChecker([]byte{0x00, 0x01})
	assert.NotNil(t, c)
	assert.Equal(t, []byte{0x00, 0x01}, c.data)
}

func TestChecker_Run(t *testing.T) {
	// todo
}

func TestChecker_Run_LoadError(t *testing.T) {
	c := NewChecker([]byte{0x00, 0x01})

	err := c.Run(CheckerOptions{})
	assert.EqualError(t, err, "error converting YAML to JSON: yaml: control characters are not allowed")
}
