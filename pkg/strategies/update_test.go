package strategies

import (
	"testing"

	"github.com/edaniszewski/chart-releaser/internal/testutils"
	version "github.com/edaniszewski/chart-releaser/pkg/semver"
	"github.com/stretchr/testify/assert"
)

func TestListStrategies(t *testing.T) {
	strategies := ListUpdateStrategies()
	assert.Len(t, strategies, 4)
	assert.Equal(t, UpdateMajor, strategies[0])
	assert.Equal(t, UpdateMinor, strategies[1])
	assert.Equal(t, UpdatePatch, strategies[2])
	assert.Equal(t, UpdateDefault, strategies[3])
}

func TestUpdateCtx_IsComplete(t *testing.T) {
	ver := version.Semver{Major: 1, Minor: 0, Patch: 0}

	tests := []struct {
		name     string
		ctx      UpdateCtx
		complete bool
	}{
		{
			name:     "missing old app version",
			ctx:      UpdateCtx{NewAppVersion: &ver, OldChartVersion: &ver},
			complete: false,
		},
		{
			name:     "missing new app version",
			ctx:      UpdateCtx{OldAppVersion: &ver, OldChartVersion: &ver},
			complete: false,
		},
		{
			name:     "missing old chart version",
			ctx:      UpdateCtx{OldAppVersion: &ver, NewAppVersion: &ver},
			complete: false,
		},
		{
			name:     "is complete",
			ctx:      UpdateCtx{OldAppVersion: &ver, NewAppVersion: &ver, OldChartVersion: &ver},
			complete: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.complete, test.ctx.IsComplete())
		})
	}
}

func TestUpdateStrategyFromString(t *testing.T) {
	tests := []struct {
		str      string
		strategy UpdateStrategy
	}{
		{str: "major", strategy: UpdateMajor},
		{str: "Major", strategy: UpdateMajor},
		{str: "MAJOR", strategy: UpdateMajor},
		{str: "minor", strategy: UpdateMinor},
		{str: "Minor", strategy: UpdateMinor},
		{str: "MINOR", strategy: UpdateMinor},
		{str: "patch", strategy: UpdatePatch},
		{str: "Patch", strategy: UpdatePatch},
		{str: "PATCH", strategy: UpdatePatch},
		{str: "default", strategy: UpdateDefault},
		{str: "Default", strategy: UpdateDefault},
		{str: "DEFAULT", strategy: UpdateDefault},
	}
	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			strategy, err := UpdateStrategyFromString(test.str)
			assert.Nil(t, err)
			assert.Equal(t, test.strategy, strategy)
		})
	}
}

func TestUpdateStrategyFromString_Error(t *testing.T) {
	strategy, err := UpdateStrategyFromString("not-a-strategy")
	assert.Error(t, err)
	assert.Empty(t, strategy)
}

func TestUpdateRelease_CtxNotComplete(t *testing.T) {
	c := UpdateCtx{
		OldChartVersion: &version.Semver{},
		NewAppVersion:   &version.Semver{},
		OldAppVersion:   &version.Semver{},
	}
	assert.True(t, c.IsComplete())
}

func TestUpdateRelease_CtxNotCompleteNoOldChart(t *testing.T) {
	c := UpdateCtx{
		NewAppVersion: &version.Semver{},
		OldAppVersion: &version.Semver{},
	}
	assert.False(t, c.IsComplete())
}

func TestUpdateRelease_CtxNotCompleteNoNewApp(t *testing.T) {
	c := UpdateCtx{
		OldChartVersion: &version.Semver{},
		OldAppVersion:   &version.Semver{},
	}
	assert.False(t, c.IsComplete())
}

func TestUpdateRelease_CtxNotCompleteNoOldApp(t *testing.T) {
	c := UpdateCtx{
		OldChartVersion: &version.Semver{},
		NewAppVersion:   &version.Semver{},
	}
	assert.False(t, c.IsComplete())
}

func TestUpdateRelease_Major(t *testing.T) {
	c := UpdateCtx{
		OldChartVersion: testutils.NewSemverP(t, "1.2.3"),
		NewAppVersion:   testutils.NewSemverP(t, "0.2.0"),
		OldAppVersion:   testutils.NewSemverP(t, "0.1.0"),
		Strategy:        UpdateMajor,
	}

	v, err := UpdateRelease(&c)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), v.Major)
	assert.Equal(t, uint64(0), v.Minor)
	assert.Equal(t, uint64(0), v.Patch)
	assert.Equal(t, "", v.Prerelease)
}

func TestUpdateRelease_Minor(t *testing.T) {
	c := UpdateCtx{
		OldChartVersion: testutils.NewSemverP(t, "1.2.3"),
		NewAppVersion:   testutils.NewSemverP(t, "0.2.0"),
		OldAppVersion:   testutils.NewSemverP(t, "0.1.0"),
		Strategy:        UpdateMinor,
	}

	v, err := UpdateRelease(&c)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), v.Major)
	assert.Equal(t, uint64(3), v.Minor)
	assert.Equal(t, uint64(0), v.Patch)
	assert.Equal(t, "", v.Prerelease)
}

func TestUpdateRelease_Patch(t *testing.T) {
	c := UpdateCtx{
		OldChartVersion: testutils.NewSemverP(t, "1.2.3"),
		NewAppVersion:   testutils.NewSemverP(t, "0.2.0"),
		OldAppVersion:   testutils.NewSemverP(t, "0.1.0"),
		Strategy:        UpdatePatch,
	}

	v, err := UpdateRelease(&c)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), v.Major)
	assert.Equal(t, uint64(2), v.Minor)
	assert.Equal(t, uint64(4), v.Patch)
	assert.Equal(t, "", v.Prerelease)
}

func TestUpdateRelease_Default(t *testing.T) {
	c := UpdateCtx{
		OldChartVersion: testutils.NewSemverP(t, "1.2.3-alpha.1"),
		NewAppVersion:   testutils.NewSemverP(t, "0.2.0-pre.1"),
		OldAppVersion:   testutils.NewSemverP(t, "0.1.0"),
		Strategy:        UpdateDefault,
	}

	v, err := UpdateRelease(&c)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), v.Major)
	assert.Equal(t, uint64(2), v.Minor)
	assert.Equal(t, uint64(3), v.Patch)
	assert.Equal(t, "alpha.2", v.Prerelease)
}

func TestUpdateRelease_ErrorNotComplete(t *testing.T) {
	c := UpdateCtx{}

	_, err := UpdateRelease(&c)
	assert.Equal(t, ErrIncompleteUpdateCtx, err)
}

func TestUpdateRelease_ErrorUnsupportedStrategy(t *testing.T) {
	c := UpdateCtx{
		OldChartVersion: &version.Semver{},
		NewAppVersion:   &version.Semver{},
		OldAppVersion:   &version.Semver{},
		Strategy:        UpdateStrategy("invalid-strat"),
	}

	_, err := UpdateRelease(&c)
	assert.EqualError(t, err, "unsupported release update strategy: invalid-strat")
}

func TestUpdateMajor(t *testing.T) {
	tests := []struct {
		name     string
		ctx      UpdateCtx
		expected version.Semver
	}{
		{
			name: "major app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "minor app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "patch app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.0-alpha.2"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0-alpha.1"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "major and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "minor and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "patch and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "major and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "minor and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "patch and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "major, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "minor, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
		{
			name: "patch, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "1.0.0"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := updateMajor(&test.ctx)
			assert.NoError(t, err)
			assert.True(t, actual.Equals(&test.expected))
		})
	}
}

func TestUpdateMajor_Error(t *testing.T) {
	tests := []struct {
		name string
		ctx  UpdateCtx
		err  error
	}{
		{
			name: "no drift: same version",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.2.3"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.3"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrNoDrift,
		},
		{
			name: "no drift: same version (with build)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.2.3+build.2"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.3+build.1"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrNoDrift,
		},
		{
			name: "new version less than old (major)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "2.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (minor)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (patch)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.3"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (prerelease)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (prerelease 2)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.2"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (build)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			// build info not used when comparing, so these should be equal
			err: version.ErrNoDrift,
		},
		{
			name: "new version less than old (build 2)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0+2"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			// build info not used when comparing, so these should be equal
			err: version.ErrNoDrift,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := updateMajor(&test.ctx)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestUpdateMinor(t *testing.T) {
	tests := []struct {
		name     string
		ctx      UpdateCtx
		expected version.Semver
	}{
		{
			name: "major app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "minor app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "patch app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.0-alpha.2"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0-alpha.1"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "major and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "minor and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "patch and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "major and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "minor and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "patch and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "major, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "minor, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
		{
			name: "patch, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.1.0"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := updateMinor(&test.ctx)
			assert.NoError(t, err)
			assert.True(t, actual.Equals(&test.expected))
		})
	}
}

func TestUpdateMinor_Error(t *testing.T) {
	tests := []struct {
		name string
		ctx  UpdateCtx
		err  error
	}{
		{
			name: "no drift: same version",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.2.3"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.3"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrNoDrift,
		},
		{
			name: "no drift: same version (with build)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.2.3+build.2"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.3+build.1"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrNoDrift,
		},
		{
			name: "new version less than old (major)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "2.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (minor)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (patch)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.3"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (prerelease)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (prerelease 2)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.2"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (build)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			// build info not used when comparing, so these should be equal
			err: version.ErrNoDrift,
		},
		{
			name: "new version less than old (build 2)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0+2"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			// build info not used when comparing, so these should be equal
			err: version.ErrNoDrift,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := updateMinor(&test.ctx)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestUpdatePatch(t *testing.T) {
	tests := []struct {
		name     string
		ctx      UpdateCtx
		expected version.Semver
	}{
		{
			name: "major app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "minor app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "patch app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.0-alpha.2"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0-alpha.1"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "major and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "minor and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "patch and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "major and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "minor and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "patch and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "major, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "minor, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "patch, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := updatePatch(&test.ctx)
			assert.NoError(t, err)
			assert.True(t, actual.Equals(&test.expected))
		})
	}
}

func TestUpdatePatch_Error(t *testing.T) {
	tests := []struct {
		name string
		ctx  UpdateCtx
		err  error
	}{
		{
			name: "no drift: same version",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.2.3"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.3"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrNoDrift,
		},
		{
			name: "no drift: same version (with build)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.2.3+build.2"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.3+build.1"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrNoDrift,
		},
		{
			name: "new version less than old (major)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "2.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (minor)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (patch)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.3"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (prerelease)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (prerelease 2)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.2"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (build)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			// build info not used when comparing, so these should be equal
			err: version.ErrNoDrift,
		},
		{
			name: "new version less than old (build 2)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0+2"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			// build info not used when comparing, so these should be equal
			err: version.ErrNoDrift,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := updatePatch(&test.ctx)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestUpdateDefault(t *testing.T) {
	tests := []struct {
		name     string
		ctx      UpdateCtx
		expected version.Semver
	}{
		{
			name: "major app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "minor app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "patch app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.0-alpha.2"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0-alpha.1"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1-pre.1"),
		},
		{
			name: "major and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1-pre.1"),
		},
		{
			name: "minor and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1-pre.1"),
		},
		{
			name: "patch and prerelease app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1-pre.1"),
		},
		{
			name: "major and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "minor and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "patch and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
		{
			name: "major, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1-pre.1"),
		},
		{
			name: "minor, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.1.0-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1-pre.1"),
		},
		{
			name: "patch, prerelease, and build app version bump",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1-pre.1"),
		},
		{
			name: "stable to prerelease",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			expected: testutils.NewSemver(t, "0.0.1-pre.1"),
		},
		{
			name: "prerelease to prerelease",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.1-alpha.1+1"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.1-alpha.1"),
			},
			expected: testutils.NewSemver(t, "0.0.1-alpha.2"),
		},
		{
			name: "prerelease to stable",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "0.0.3"),
				OldAppVersion:   testutils.NewSemverP(t, "0.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.1-alpha.4"),
			},
			expected: testutils.NewSemver(t, "0.0.1"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := updateDefault(&test.ctx)
			assert.NoError(t, err)
			assert.Equal(t, test.expected.String(), actual.String())
		})
	}
}

func TestUpdateDefault_Error(t *testing.T) {
	tests := []struct {
		name string
		ctx  UpdateCtx
		err  error
	}{
		{
			name: "no drift: same version",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.2.3"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.3"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrNoDrift,
		},
		{
			name: "no drift: same version (with build)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.2.3+build.2"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.3+build.1"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrNoDrift,
		},
		{
			name: "new version less than old (major)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "2.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (minor)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "1.2.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (patch)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.3"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (prerelease)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (prerelease 2)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0-alpha.2"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			err: version.ErrVersionLessThanBase,
		},
		{
			name: "new version less than old (build)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			// build info not used when comparing, so these should be equal
			err: version.ErrNoDrift,
		},
		{
			name: "new version less than old (build 2)",
			ctx: UpdateCtx{
				NewAppVersion:   testutils.NewSemverP(t, "1.0.0+1"),
				OldAppVersion:   testutils.NewSemverP(t, "1.0.0+2"),
				OldChartVersion: testutils.NewSemverP(t, "0.0.0"),
			},
			// build info not used when comparing, so these should be equal
			err: version.ErrNoDrift,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := updateDefault(&test.ctx)
			assert.Equal(t, test.err, err)
		})
	}
}
