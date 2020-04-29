package utils

import (
	"fmt"
	"strings"

	context "github.com/edaniszewski/chart-releaser/pkg/v1/ctx"
)

// ParseRepository parses a full repository string into its constituent parts.
// If the repository identifier is missing a "type" (e.g. github.com), it is
// given the default value of github.
func ParseRepository(repo string) (context.Repository, error) {
	var repository context.Repository

	parts := strings.Split(repo, "/")
	if len(parts) == 2 {
		// If there are only two parts, assume the default type of GitHub.
		repository.Owner = parts[0]
		repository.Name = parts[1]
		repository.Type = context.RepoGithub

	} else if len(parts) == 3 {
		// Check that the first part matches a supported repository type.
		rt, err := context.RepoTypeFromString(parts[0])
		if err != nil {
			return repository, err
		}
		repository.Owner = parts[1]
		repository.Name = parts[2]
		repository.Type = rt

	} else {
		// Unexpected number of parts.
		return repository, fmt.Errorf("unexpected repository string format - should be in the form of REPO/OWNER/NAME")
	}
	return repository, nil
}
