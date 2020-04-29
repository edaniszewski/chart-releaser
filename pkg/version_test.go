package pkg

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersionInfo(t *testing.T) {
	info := NewVersionInfo()

	// The version info here should be empty since their corresponding global
	// variables were not set, e.g. via build-time args.
	assert.Equal(t, "", info.Version)
	assert.Equal(t, "", info.Commit)
	assert.Equal(t, "", info.Tag)
	assert.Equal(t, "", info.GoVersion)
	assert.Equal(t, "", info.BuildDate)
	assert.Equal(t, runtime.Compiler, info.Compiler)
	assert.Equal(t, runtime.GOOS, info.OS)
	assert.Equal(t, runtime.GOARCH, info.Arch)
}
