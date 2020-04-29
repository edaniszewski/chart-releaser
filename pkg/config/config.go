package config

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// DefaultFile is the default configuration file name for chart-releaser.
const DefaultFile = ".chartreleaser.yml"

// Errors related to loading and parsing chart-releaser configuration.
var (
	ErrConfigExists    = errors.New(".chartreleaser,yml config file already exists")
	ErrNoConfig        = errors.New(".chartreleaser.yml file not found")
	ErrNoConfigVersion = errors.New("no version specified in config")
)

// VersionedConfig is a lightweight struct which is used as a first step
// in loading config files as a means to get the config version. Once the
// config version is known, the correct versioned struct can be used to
// unmarshal the data.
type VersionedConfig struct {
	Version *string `yaml:"version,omitempty"`

	data []byte
	path string
}

// GetData gets the data, as bytes, for the config file.
func (c *VersionedConfig) GetData() []byte {
	if c == nil {
		return []byte(nil)
	}
	return c.data
}

// GetVersion gets the version scheme for the config file.
func (c *VersionedConfig) GetVersion() string {
	if c == nil {
		return ""
	}
	if c.Version == nil {
		return ""
	}
	return *c.Version
}

// GetPath gets the path which the config was loaded from.
func (c *VersionedConfig) GetPath() string {
	if c == nil {
		return ""
	}
	return c.path
}

// Load the config file at the specified path into a VersionedConfig.
func Load(p string) (*VersionedConfig, error) {
	path, err := GetConfigPath(p)
	if err != nil {
		return nil, err
	}

	// Verify that the config path exists.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, ErrNoConfig
	}

	// Read from the file and figure out what the config version is.
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var v VersionedConfig
	if err := yaml.Unmarshal(contents, &v); err != nil {
		return nil, err
	}
	if v.Version == nil {
		return nil, ErrNoConfigVersion
	}
	v.data = contents
	v.path = path

	return &v, nil
}
