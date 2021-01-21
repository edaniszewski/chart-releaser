package version

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newSemver(t *testing.T, v string) Semver {
	sv, err := Load(v)
	if err != nil {
		t.Fatal(err)
	}
	return sv
}

func TestLoad(t *testing.T) {
	tests := []struct {
		v        string
		expected Semver
	}{
		{
			v: "0.0.0",
			expected: Semver{
				Major:      0,
				Minor:      0,
				Patch:      0,
				Prerelease: "",
				Build:      "",
			},
		},
		{
			v: "v0.0.0",
			expected: Semver{
				Major:      0,
				Minor:      0,
				Patch:      0,
				Prerelease: "",
				Build:      "",
				hasPrefix:  true,
			},
		},
		{
			v: "0.0.1",
			expected: Semver{
				Major:      0,
				Minor:      0,
				Patch:      1,
				Prerelease: "",
				Build:      "",
			},
		},
		{
			v: "v0.0.1",
			expected: Semver{
				Major:      0,
				Minor:      0,
				Patch:      1,
				Prerelease: "",
				Build:      "",
				hasPrefix:  true,
			},
		},
		{
			v: "1.2.3",
			expected: Semver{
				Major:      1,
				Minor:      2,
				Patch:      3,
				Prerelease: "",
				Build:      "",
			},
		},
		{
			v: "v1.2.3",
			expected: Semver{
				Major:      1,
				Minor:      2,
				Patch:      3,
				Prerelease: "",
				Build:      "",
				hasPrefix:  true,
			},
		},
		{
			v: "4.2.0-alpha",
			expected: Semver{
				Major:      4,
				Minor:      2,
				Patch:      0,
				Prerelease: "alpha",
				Build:      "",
			},
		},
		{
			v: "v4.2.0-alpha",
			expected: Semver{
				Major:      4,
				Minor:      2,
				Patch:      0,
				Prerelease: "alpha",
				Build:      "",
				hasPrefix:  true,
			},
		},
		{
			v: "2.2.2-rc.9",
			expected: Semver{
				Major:      2,
				Minor:      2,
				Patch:      2,
				Prerelease: "rc.9",
				Build:      "",
			},
		},
		{
			v: "v2.2.2-rc.9",
			expected: Semver{
				Major:      2,
				Minor:      2,
				Patch:      2,
				Prerelease: "rc.9",
				Build:      "",
				hasPrefix:  true,
			},
		},
		{
			v: "0.0.1+build2",
			expected: Semver{
				Major:      0,
				Minor:      0,
				Patch:      1,
				Prerelease: "",
				Build:      "build2",
			},
		},
		{
			v: "v0.0.1+build2",
			expected: Semver{
				Major:      0,
				Minor:      0,
				Patch:      1,
				Prerelease: "",
				Build:      "build2",
				hasPrefix:  true,
			},
		},
		{
			v: "0.2.0+build.3",
			expected: Semver{
				Major:      0,
				Minor:      2,
				Patch:      0,
				Prerelease: "",
				Build:      "build.3",
			},
		},
		{
			v: "v0.2.0+build.3",
			expected: Semver{
				Major:      0,
				Minor:      2,
				Patch:      0,
				Prerelease: "",
				Build:      "build.3",
				hasPrefix:  true,
			},
		},
		{
			v: "15.10.200-beta.9+build10.1",
			expected: Semver{
				Major:      15,
				Minor:      10,
				Patch:      200,
				Prerelease: "beta.9",
				Build:      "build10.1",
			},
		},
		{
			v: "v15.10.200-beta.9+build10.1",
			expected: Semver{
				Major:      15,
				Minor:      10,
				Patch:      200,
				Prerelease: "beta.9",
				Build:      "build10.1",
				hasPrefix:  true,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.v, func(t *testing.T) {
			actual, err := Load(test.v)
			assert.NoError(t, err)
			assert.True(t, test.expected.Equals(&actual))
		})
	}
}

func TestLoad_ParseError(t *testing.T) {
	_, err := Load("invalid version string")
	assert.EqualError(t, err, "No Major.Minor.Patch elements found")
}

func TestLoad_PRVersionError(t *testing.T) {
	_, err := Load("1.2.3-000")
	assert.EqualError(t, err, "Numeric PreRelease version must not contain leading zeroes \"000\"")
}

func TestSemver_String(t *testing.T) {
	s := Semver{
		Major:      1,
		Minor:      2,
		Patch:      3,
		Prerelease: "alpha.1",
		Build:      "build2",
	}

	assert.Equal(t, "1.2.3-alpha.1+build2", s.String())
}

func TestSemver_StringWithVersionPrefix(t *testing.T) {
	s := Semver{
		Major:      1,
		Minor:      2,
		Patch:      3,
		Prerelease: "alpha.1",
		Build:      "build2",
		hasPrefix:  true,
	}

	assert.Equal(t, "v1.2.3-alpha.1+build2", s.String())
}

func TestSemver_EqualsEmpty(t *testing.T) {
	v1 := Semver{}
	v2 := Semver{}

	assert.True(t, v1.Equals(&v2))
}

func TestSemver_EqualsTrue(t *testing.T) {
	v1 := Semver{Major: 1, Minor: 4, Patch: 2, Prerelease: "alpha"}
	v2 := Semver{Major: 1, Minor: 4, Patch: 2, Prerelease: "alpha"}

	assert.True(t, v1.Equals(&v2))
}

func TestSemver_EqualsFalse(t *testing.T) {
	v1 := Semver{Major: 1, Minor: 3, Patch: 5}
	v2 := Semver{Major: 1, Minor: 2, Patch: 5}

	assert.False(t, v1.Equals(&v2))
}

func TestSemver_WithVersionContext(t *testing.T) {
	v := Semver{Major: 2}
	assert.Nil(t, v.versionCtx)

	v2 := Semver{Major: 1}
	v.WithVersionContext(&v2)

	assert.Equal(t, &v2, v.versionCtx)
}

func TestSemver_Copy(t *testing.T) {
	v1 := Semver{
		Major:      1,
		Minor:      2,
		Patch:      3,
		Prerelease: "alpha",
		Build:      "build1",
	}

	assert.True(t, v1.Equals(v1.Copy()))
}

func TestSemver_Compare(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{
			name:     "%s equals %s",
			v1:       "0.0.0",
			v2:       "0.0.0",
			expected: 0,
		},
		{
			name:     "%s equals %s",
			v1:       "1.0.2",
			v2:       "1.0.2",
			expected: 0,
		},
		{
			name:     "%s equals %s",
			v1:       "1.4.0-alpha1",
			v2:       "1.4.0-alpha1",
			expected: 0,
		},
		{
			name:     "%s equals %s",
			v1:       "0.4.0+build4",
			v2:       "0.4.0+build4",
			expected: 0,
		},
		{
			name:     "%s equals %s",
			v1:       "1.30.1-alpha.3+build.1",
			v2:       "1.30.1-alpha.3+build.1",
			expected: 0,
		},

		{
			name:     "%s less than %s",
			v1:       "0.3.1",
			v2:       "1.0.2",
			expected: -1,
		},
		{
			name:     "%s less than %s",
			v1:       "1.4.0-alpha0",
			v2:       "1.4.0-alpha1",
			expected: -1,
		},
		{
			name: "%s equal to %s",
			v1:   "0.4.0+build2",
			v2:   "0.4.0+build4",
			// per the semver spec, build level is not used during
			// version comparison
			expected: 0,
		},
		{
			name:     "%s less than %s",
			v1:       "1.30.1-alpha.3+build.10",
			v2:       "1.30.1-alpha.4+build.1",
			expected: -1,
		},
		{
			name:     "%s less than %s",
			v1:       "1.30.1-alpha.2",
			v2:       "1.30.1",
			expected: -1,
		},

		{
			name:     "%s greater than %s",
			v1:       "1.1.1",
			v2:       "1.0.2",
			expected: 1,
		},
		{
			name:     "%s greater than %s",
			v1:       "1.4.0-alpha4",
			v2:       "1.4.0-alpha1",
			expected: 1,
		},
		{
			name: "%s equal to %s",
			v1:   "0.4.0+build4",
			v2:   "0.4.0+build1",
			// per the semver spec, build level is not used during
			// version comparison
			expected: 0,
		},
		{
			name:     "%s greater than %s",
			v1:       "1.30.1-alpha.5+build.1",
			v2:       "1.30.1-alpha.3+build.10",
			expected: 1,
		},
		{
			name:     "%s greater than %s",
			v1:       "1.30.1",
			v2:       "1.30.1-alpha.1",
			expected: 1,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf(test.name, test.v1, test.v2), func(t *testing.T) {
			ver1, err := Load(test.v1)
			assert.NoError(t, err)

			ver2, err := Load(test.v2)
			assert.NoError(t, err)

			assert.Equal(t, test.expected, ver1.Compare(&ver2))
		})
	}
}

func TestSemver_IncrementNew(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{
			level:    LevelMajor,
			expected: "1.0.0",
		},
		{
			level:    LevelMinor,
			expected: "0.2.0",
		},
		{
			level:    LevelPatch,
			expected: "0.1.3",
		},
		{
			level:    LevelPrerelease,
			expected: "0.1.2-pre.1",
		},
		{
			level:    LevelNone,
			expected: "0.1.2",
		},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			v, err := Load("0.1.2")
			assert.NoError(t, err)

			v2 := v.IncrementNew(test.level)
			assert.Equal(t, test.expected, v2.String())
		})
	}
}

func TestSemver_Increment(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{
			level:    LevelMajor,
			expected: "1.0.0",
		},
		{
			level:    LevelMinor,
			expected: "0.2.0",
		},
		{
			level:    LevelPatch,
			expected: "0.1.3",
		},
		{
			level:    LevelPrerelease,
			expected: "0.1.2-pre.1",
		},
		{
			level:    LevelNone,
			expected: "0.1.2",
		},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			v, err := Load("0.1.2")
			assert.NoError(t, err)

			v.Increment(test.level)
			assert.Equal(t, test.expected, v.String())
		})
	}
}

func TestSemver_IncrementUnsupportedLevel(t *testing.T) {
	assert.Panics(t, func() {
		v, err := Load("0.1.2")
		assert.NoError(t, err)

		v.Increment(Level(100))
	})
}

func TestSemver_FindDrift(t *testing.T) {
	baseVersion1 := newSemver(t, "0.1.0-alpha.1")

	tests := []struct {
		name       string
		expected   Level
		prerelease bool
		v          Semver
		base       *Semver
	}{
		{
			name:       "major drift",
			expected:   LevelMajor,
			prerelease: false,
			v:          newSemver(t, "1.0.0"),
		}, {
			name:       "major with prerelease drift",
			expected:   LevelMajor,
			prerelease: true,
			v:          newSemver(t, "1.0.0-alpha.1"),
		},
		{
			name:       "minor drift",
			expected:   LevelMinor,
			prerelease: false,
			v:          newSemver(t, "0.2.0"),
		},
		{
			name:       "minor with prerelease drift",
			expected:   LevelMinor,
			prerelease: true,
			v:          newSemver(t, "0.2.0-alpha.1"),
		},
		{
			name:       "patch drift",
			expected:   LevelPatch,
			prerelease: false,
			v:          newSemver(t, "0.1.1"),
		},
		{
			name:       "patch with prerelease drift",
			expected:   LevelPatch,
			prerelease: true,
			v:          newSemver(t, "0.1.1-alpha.1"),
		},
		{
			name:       "prerelease drift",
			expected:   LevelPrerelease,
			prerelease: true,
			v:          newSemver(t, "0.1.0-alpha.2"),
			base:       &baseVersion1,
		},
		{
			name:       "major and minor drift",
			expected:   LevelMajor,
			prerelease: false,
			v:          newSemver(t, "1.2.0"),
		},
		{
			name:       "minor and patch drift",
			expected:   LevelMinor,
			prerelease: false,
			v:          newSemver(t, "0.2.1"),
		},
		{
			name:       "patch and prerelease drift",
			expected:   LevelPatch,
			prerelease: true,
			v:          newSemver(t, "0.1.1-alpha"),
		},
		{
			name:       "prerelease and build drift",
			expected:   LevelPrerelease,
			prerelease: true,
			v:          newSemver(t, "0.1.0-alpha.2+1"),
			base:       &baseVersion1,
		},
	}

	defaultBase := newSemver(t, "0.1.0")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			base := &defaultBase
			if test.base != nil {
				base = test.base
			}

			actual, pre, err := test.v.FindDrift(base)
			assert.NoError(t, err)
			assert.Equal(t, test.prerelease, pre)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestSemver_FindDrift_ErrorLessThanBase(t *testing.T) {
	ver := newSemver(t, "0.1.0")
	base := newSemver(t, "0.2.0")

	v, pre, err := ver.FindDrift(&base)
	assert.Equal(t, ErrVersionLessThanBase, err)
	assert.False(t, pre)
	assert.Equal(t, LevelNone, v)
}

func TestSemver_FindDrift_ErrorNoDrift(t *testing.T) {
	ver := newSemver(t, "0.1.0")
	base := newSemver(t, "0.1.0")

	v, pre, err := ver.FindDrift(&base)
	assert.Equal(t, ErrNoDrift, err)
	assert.False(t, pre)
	assert.Equal(t, LevelNone, v)
}

func TestSemver_IncrementMajor(t *testing.T) {
	v := newSemver(t, "1.2.3")

	v.IncrementMajor()
	assert.Equal(t, uint64(2), v.Major)
	assert.Equal(t, uint64(0), v.Minor)
	assert.Equal(t, uint64(0), v.Patch)
	assert.Equal(t, "", v.Prerelease)
}

func TestSemver_IncrementMinor(t *testing.T) {
	v := newSemver(t, "1.2.3")

	v.IncrementMinor()
	assert.Equal(t, uint64(1), v.Major)
	assert.Equal(t, uint64(3), v.Minor)
	assert.Equal(t, uint64(0), v.Patch)
	assert.Equal(t, "", v.Prerelease)
}

func TestSemver_IncrementPatch(t *testing.T) {
	v := newSemver(t, "1.2.3")

	v.IncrementPatch()
	assert.Equal(t, uint64(1), v.Major)
	assert.Equal(t, uint64(2), v.Minor)
	assert.Equal(t, uint64(4), v.Patch)
	assert.Equal(t, "", v.Prerelease)
}

func TestSemver_IncrementPrerelease_New(t *testing.T) {
	v := newSemver(t, "1.2.3")

	v.IncrementPrerelease()
	assert.Equal(t, uint64(1), v.Major)
	assert.Equal(t, uint64(2), v.Minor)
	assert.Equal(t, uint64(3), v.Patch)
	assert.Equal(t, "pre.1", v.Prerelease)
}

func TestSemver_IncrementPrerelease_ExistingWithNumeric(t *testing.T) {
	v := newSemver(t, "1.2.3-alpha.1")

	v.IncrementPrerelease()
	assert.Equal(t, uint64(1), v.Major)
	assert.Equal(t, uint64(2), v.Minor)
	assert.Equal(t, uint64(3), v.Patch)
	assert.Equal(t, "alpha.2", v.Prerelease)
}

func TestSemver_IncrementPrerelease_ExistingWithNumeric2(t *testing.T) {
	v := newSemver(t, "1.2.3-alpha.1.2")

	v.IncrementPrerelease()
	assert.Equal(t, uint64(1), v.Major)
	assert.Equal(t, uint64(2), v.Minor)
	assert.Equal(t, uint64(3), v.Patch)
	assert.Equal(t, "alpha.1.3", v.Prerelease)
}

func TestSemver_IncrementPrerelease_ExistingWithNumeric3(t *testing.T) {
	v := newSemver(t, "1.2.3-alpha.1.2.test")

	v.IncrementPrerelease()
	assert.Equal(t, uint64(1), v.Major)
	assert.Equal(t, uint64(2), v.Minor)
	assert.Equal(t, uint64(3), v.Patch)
	assert.Equal(t, "alpha.1.3.test", v.Prerelease)
}

func TestSemver_IncrementPrerelease_ExistingWithoutNumeric(t *testing.T) {
	v := newSemver(t, "1.2.3-alpha")

	v.IncrementPrerelease()
	assert.Equal(t, uint64(1), v.Major)
	assert.Equal(t, uint64(2), v.Minor)
	assert.Equal(t, uint64(3), v.Patch)
	assert.Equal(t, "alpha.1", v.Prerelease)
}
