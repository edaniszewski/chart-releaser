package strategies

import (
	"errors"
	"fmt"
	"strings"

	"github.com/apex/log"
	version "github.com/edaniszewski/chart-releaser/pkg/semver"
)

// Errors relating to updating release versions.
var (
	ErrIncompleteUpdateCtx = errors.New("incomplete update context")
)

// UpdateStrategy is a strategy used to increment the semantic version of
// a Helm Chart.
type UpdateStrategy string

// The update strategies that chart-releaser supports.
const (
	UpdateMajor   UpdateStrategy = "major"
	UpdateMinor   UpdateStrategy = "minor"
	UpdatePatch   UpdateStrategy = "patch"
	UpdateDefault UpdateStrategy = "default"
)

// ListUpdateStrategies returns a list of all supported UpdateStrategies.
func ListUpdateStrategies() []UpdateStrategy {
	return []UpdateStrategy{
		UpdateMajor,
		UpdateMinor,
		UpdatePatch,
		UpdateDefault,
	}
}

// UpdateCtx contains information needed to perform an update. Since the
// update to the Chart version may depend on the level at which the application
// version was incremented, this needs to hold both the new and old app version in
// order to detect drift, as well as the old chart version, which will be the seed
// for the new chart version.
type UpdateCtx struct {
	OldAppVersion   *version.Semver
	NewAppVersion   *version.Semver
	OldChartVersion *version.Semver

	Strategy UpdateStrategy
}

// IsComplete checks whether all versions needed by the UpdateCtx are specified.
func (ctx *UpdateCtx) IsComplete() bool {
	return !(ctx.OldAppVersion == nil || ctx.NewAppVersion == nil || ctx.OldChartVersion == nil)
}

// UpdateStrategyFromString returns the UpdateStrategy corresponding to the provided string.
func UpdateStrategyFromString(s string) (UpdateStrategy, error) {
	switch strings.ToLower(s) {
	case "major":
		return UpdateMajor, nil
	case "minor":
		return UpdateMinor, nil
	case "patch":
		return UpdatePatch, nil
	case "default":
		return UpdateDefault, nil
	default:
		return "", fmt.Errorf("not a valid strategy: %s", s)
	}
}

// UpdateRelease applies an update strategy to generate a new Chart version.
func UpdateRelease(ctx *UpdateCtx) (version.Semver, error) {
	if !ctx.IsComplete() {
		return version.Semver{}, ErrIncompleteUpdateCtx
	}

	switch ctx.Strategy {
	case UpdateMajor:
		return updateMajor(ctx)
	case UpdateMinor:
		return updateMinor(ctx)
	case UpdatePatch:
		return updatePatch(ctx)
	case UpdateDefault:
		return updateDefault(ctx)
	default:
		return version.Semver{}, fmt.Errorf("unsupported release update strategy: %s", ctx.Strategy)
	}
}

// updateMajor is executed upon updating a release via UpdateCtx for the "major" strategy.
//
// Under this strategy, the Chart's major version is incremented for any change to the app version.
func updateMajor(ctx *UpdateCtx) (version.Semver, error) {
	_, _, err := ctx.NewAppVersion.FindDrift(ctx.OldAppVersion)
	if err != nil {
		log.WithError(err).Error("update major: failed to find application version drift")
		return version.Semver{}, err
	}

	// The "major" strategy says to bump the major version of the chart for any
	// change to the app version.
	return ctx.OldChartVersion.IncrementNew(version.LevelMajor), nil
}

// updateMinor is executed upon updating a release via UpdateCtx for the "minor" strategy.
//
// Under this strategy, the Chart's minor version is incremented for any change to the app version.
func updateMinor(ctx *UpdateCtx) (version.Semver, error) {
	_, _, err := ctx.NewAppVersion.FindDrift(ctx.OldAppVersion)
	if err != nil {
		log.WithError(err).Error("update minor: failed to find application version drift")
		return version.Semver{}, err
	}

	// The "minor" strategy says to bump the minor version of the chart for any
	// change to the app version.
	return ctx.OldChartVersion.IncrementNew(version.LevelMinor), nil
}

// updatePatch is executed upon updating a release via UpdateCtx for the "patch" strategy.
//
// Under this strategy, the Chart's patch version is incremented for any change to the app version.
func updatePatch(ctx *UpdateCtx) (version.Semver, error) {
	_, _, err := ctx.NewAppVersion.FindDrift(ctx.OldAppVersion)
	if err != nil {
		log.WithError(err).Error("update patch: failed to find application version drift")
		return version.Semver{}, err
	}

	// The "patch" strategy says to bump the patch version of the chart for any
	// change to the app version.
	return ctx.OldChartVersion.IncrementNew(version.LevelPatch), nil
}

// updateDefault is executed upon updating a release via UpdateCtx for the "default" strategy.
//
// Under this strategy, the Chart's patch version is incremented for any change to the app version,
// except in the case where pre-release info is contained within the app version. In such case, the
// Chart version either has a prerelease assigned to it, or has an existing prerelease bumped.
func updateDefault(ctx *UpdateCtx) (version.Semver, error) {
	_, prerelease, err := ctx.NewAppVersion.FindDrift(ctx.OldAppVersion)
	if err != nil {
		log.WithError(err).Error("update default: failed to find application version drift")
		return version.Semver{}, err
	}

	var newVersion version.Semver
	if !prerelease {
		// If the drift does not have a prerelease, we are either incrementing from
		// a non-prerelease version (e.g. 0.1.0 -> 0.1.1) or from  a prerelease version
		// which has become stable (e.g. 0.1.0-pre.1 -> 0.1.0). If the old chart version
		// has a prerelease assigned to it, simply remove the prerelease. Otherwise increment
		// the version.
		if ctx.OldChartVersion.Prerelease != "" {
			newVersion = *ctx.OldChartVersion.Copy()
			newVersion.Prerelease = ""
		} else {
			newVersion = ctx.OldChartVersion.IncrementNew(version.LevelPatch)
		}
	} else {
		if ctx.OldChartVersion.Prerelease == "" {
			newVersion = ctx.OldChartVersion.IncrementNew(version.LevelPatch)
			newVersion.IncrementPrerelease()
		} else {
			newVersion = *ctx.OldChartVersion.Copy()
			newVersion.IncrementPrerelease()
		}
	}
	return newVersion, nil
}
