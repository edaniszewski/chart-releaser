package pkg

import "runtime"

// Build-time arguments are used to populate these variables to generate
// build-specific version information.
var (
	Version   string
	Commit    string
	Tag       string
	GoVersion string
	BuildDate string
)

// VersionInfo contains runtime and build time information about the application.
type VersionInfo struct {
	Version   string
	Commit    string
	Tag       string
	GoVersion string
	BuildDate string
	Compiler  string
	OS        string
	Arch      string
}

// NewVersionInfo creates a struct which holds all of the runtime and
// build-time supplied variables describing the version and build state
// for the application.
func NewVersionInfo() VersionInfo {
	return VersionInfo{
		Version:   Version,
		Commit:    Commit,
		Tag:       Tag,
		GoVersion: GoVersion,
		BuildDate: BuildDate,
		Compiler:  runtime.Compiler,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// LatestVersion contains information about the latest version of the project
// and the current version of the project.
type LatestVersion struct {
	Latest    string
	Installed string
	Status    string
}
