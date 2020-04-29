package version

import (
	"errors"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/blang/semver"
)

// Errors for operations on semantic versions.
var (
	ErrVersionLessThanBase  = errors.New("version is less than the starting base version")
	ErrNoDrift              = errors.New("no version drift detected, versions are equal")
	ErrFailedVersionCompare = errors.New("failed to compare two versions")
)

// Level is the level of a component of a semantic version.
type Level uint8

// The different level definitions for semantic version components.
const (
	LevelNone Level = iota
	LevelPrerelease
	LevelPatch
	LevelMinor
	LevelMajor
)

// Semver represents a parsed semantic version.
type Semver struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	Prerelease string
	Build      string

	prerelease []semver.PRVersion
	build      []semver.PRVersion
	versionCtx *Semver
	v          semver.Version
}

// Load a new Semver from a semantic version string.
func Load(version string) (Semver, error) {
	v, err := semver.Parse(version)
	if err != nil {
		return Semver{}, err
	}

	var pre []string
	for _, s := range v.Pre {
		pre = append(pre, s.String())
	}

	// While not quite what was intended by the semver lib, we parse the build
	// strings as pre-release strings in order to use the utility functions defined
	// on the PRVersion type, namely to check whether it is numeric or not.
	var build []semver.PRVersion
	for _, s := range v.Build {
		v, err := semver.NewPRVersion(s)
		if err != nil {
			return Semver{}, err
		}
		build = append(build, v)
	}

	return Semver{
		Major:      v.Major,
		Minor:      v.Minor,
		Patch:      v.Patch,
		Prerelease: strings.Join(pre, "."),
		Build:      strings.Join(v.Build, "."),
		prerelease: v.Pre,
		build:      build,
		v:          v,
	}, nil
}

// String returns the string representation of the semantic version.
func (s *Semver) String() string {
	str := fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)
	if s.Prerelease != "" {
		str += "-"
		str += s.Prerelease
	}
	if s.Build != "" {
		str += "+"
		str += s.Build
	}
	return str
}

// Equals checks whether one Semver is equal to another. Equality is determined
// by comparing the major/minor/patch/prerelease/build values. Since those values
// are all encoded into the string representation of the Semver, this just compares
// those strings.
func (s *Semver) Equals(other *Semver) bool {
	return s.String() == other.String()
}

// WithVersionContext adds a contextual semantic version to the Semver.
func (s *Semver) WithVersionContext(v *Semver) *Semver {
	s.versionCtx = v
	return s
}

// Copy the Semver into a new instance.
func (s *Semver) Copy() *Semver {
	return &Semver{
		Major:      s.Major,
		Minor:      s.Minor,
		Patch:      s.Patch,
		Prerelease: s.Prerelease,
		Build:      s.Build,
		prerelease: s.prerelease,
		build:      s.build,
		versionCtx: s.versionCtx,
	}
}

// Compare two versions to determine if one is less than (-1), equal to (0),
// or greater than (1) the other.
func (s *Semver) Compare(other *Semver) int {
	return s.v.Compare(other.v)
}

// IncrementNew creates a copy of the Semver and increments the new copy at the
// specified level.
func (s *Semver) IncrementNew(level Level) Semver {
	v := s.Copy()
	v.Increment(level)
	return *v
}

// Increment the Semver at the specified level.
func (s *Semver) Increment(level Level) {
	switch level {
	case LevelMajor:
		s.IncrementMajor()
	case LevelMinor:
		s.IncrementMinor()
	case LevelPatch:
		s.IncrementPatch()
	case LevelPrerelease:
		s.IncrementPrerelease()
	case LevelNone:
		// If there is no level set, there is nothing to increment.
		return
	default:
		panic(fmt.Sprintf("attempting to increment version with unsupported level: %v", level))
	}
}

// IncrementMajor increments the major version of the Semver. In doing so, it resets
// lower version components, e.g. a major version bump from 0.1.2 would be 1.0.0.
func (s *Semver) IncrementMajor() {
	// Increment the major version.
	s.Major++

	// Reset lower version components.
	s.Minor = 0
	s.Patch = 0
	s.Prerelease = ""
	s.Build = ""

	// Recreate the underlying version model.
	v, err := semver.Parse(s.String())
	if err != nil {
		// Panic. Since the string builder should create a valid Semver,
		// an error in parsing should be deemed fatal.
		log.WithField("version", s.String()).Error("error: failed to create new Semver for major increment")
		panic(err)
	}
	s.v = v
}

// IncrementMinor increments the minor version of the Semver. In doing so, it resets
// lower version components, e.g. a minor version bump from 0.1.2 would be 0.2.0.
func (s *Semver) IncrementMinor() {
	// Increment the minor version.
	s.Minor++

	// Reset lower version components.
	s.Patch = 0
	s.Prerelease = ""
	s.Build = ""

	// Recreate the underlying version model.
	v, err := semver.Parse(s.String())
	if err != nil {
		// Panic. Since the string builder should create a valid Semver,
		// an error in parsing should be deemed fatal.
		log.WithField("version", s.String()).Error("error: failed to create new Semver for minor increment")
		panic(err)
	}
	s.v = v
}

// IncrementPatch increments the patch version of the Semver. In doing so, it resets
// lower version components, e.g. a patch version bump from 0.1.2 would be 0.1.3.
func (s *Semver) IncrementPatch() {
	// Increment the patch version.
	s.Patch++

	// Reset lower version components.
	s.Prerelease = ""
	s.Build = ""

	// Recreate the underlying version model.
	v, err := semver.Parse(s.String())
	if err != nil {
		// Panic. Since the string builder should create a valid Semver,
		// an error in parsing should be deemed fatal.
		log.WithField("version", s.String()).Error("error: failed to create new Semver for patch increment")
		panic(err)
	}
	s.v = v
}

// IncrementPrerelease increments the prerelease version of the Semver.
func (s *Semver) IncrementPrerelease() {
	if s.Prerelease == "" {
		// If there is currently no pre-release, add the default "-pre.1"
		comp1, _ := semver.NewPRVersion("pre")
		comp2, _ := semver.NewPRVersion("1")
		s.prerelease = []semver.PRVersion{comp1, comp2}
	} else {
		// Otherwise, there is an existing pre-release, so increment the
		// value there.
		var incremented bool
		for i := len(s.prerelease) - 1; i >= 0; i-- {
			if s.prerelease[i].IsNumeric() {
				s.prerelease[i].VersionNum++
				incremented = true
				break
			}
		}
		if !incremented {
			newComp, _ := semver.NewPRVersion("1")
			s.prerelease = append(s.prerelease, newComp)
		}
	}

	var pre []string
	for _, p := range s.prerelease {
		pre = append(pre, p.String())
	}
	s.Prerelease = strings.Join(pre, ".")
}

// FindDrift finds the first instance of version drift walking down the version
// components.
//
// This will return an error if the version is less than the base version it is
// being compared to, or if the versions are equal.
//
// Additionally, it tracks whether the pre-release version was incremented, allowing
// the caller to determine whether the drift is a normal version drift (0.0.1 -> 0.0.2),
// a prerelease drift (0.1.0-alpha.1 -> 0.1.0-alpha.2), or a version with prerelease drift
// (0.1.0 -> 0.2.0-alpha.1).
func (s *Semver) FindDrift(base *Semver) (Level, bool, error) {
	// Only find the drift if the current version is greater than the
	// supplied 'base' version. For our use cases, we don't care about
	// absolute drift, only relative drift.
	switch s.Compare(base) {
	case -1:
		return LevelNone, false, ErrVersionLessThanBase
	case 0:
		return LevelNone, false, ErrNoDrift
	case 1:
		var level = LevelNone
		var hasPrerelease bool

		if s.Major != base.Major {
			level = LevelMajor
		} else if s.Minor != base.Minor {
			level = LevelMinor
		} else if s.Patch != base.Patch {
			level = LevelPatch
		} else if s.Prerelease != base.Prerelease {
			level = LevelPrerelease
		}

		if s.Prerelease != "" {
			hasPrerelease = true
		}

		return level, hasPrerelease, nil

	default:
		return LevelNone, false, ErrFailedVersionCompare
	}
}
