package testutils

import (
	"testing"

	version "github.com/edaniszewski/chart-releaser/pkg/semver"
)

// NewSemver is a test utility to create a new SemVer from a string.
func NewSemver(t *testing.T, v string) version.Semver {
	sv, err := version.Load(v)
	if err != nil {
		t.Fatal(err)
	}
	return sv
}

// NewSemverP is a test utility to create a new SemVer pointer from a string.
func NewSemverP(t *testing.T, v string) *version.Semver {
	sv := NewSemver(t, v)
	return &sv
}
