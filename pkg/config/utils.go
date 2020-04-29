package config

import (
	"os"
	"path/filepath"
)

// GetConfigPath gets the path to the chart-releaser config file.
func GetConfigPath(path string) (string, error) {
	if path == "" {
		// If no path is provided, assume the current working directory.
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = filepath.Join(dir, DefaultFile)
	} else {
		// If a path is provided and it is a directory, look for the default
		// config file within that directory.
		info, err := os.Stat(path)
		if err != nil {
			return "", err
		}
		if info.IsDir() {
			path = filepath.Join(path, DefaultFile)
		}
	}
	return path, nil
}
