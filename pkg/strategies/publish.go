package strategies

import (
	"fmt"
	"strings"
)

// PublishStrategy defines the behavior that chart-releaser will use
// to publish updates and changes to a chart repository.
type PublishStrategy string

// The publish strategies supported by chart-releaser.
const (
	PublishCommit      PublishStrategy = "commit"
	PublishPullRequest PublishStrategy = "pull request"
)

// ListPublishStrategies returns a slice of all supported PublishStrategies.
func ListPublishStrategies() []PublishStrategy {
	return []PublishStrategy{
		PublishCommit,
		PublishPullRequest,
	}
}

// PublishStrategyFromString returns the PublishStrategy corresponding to the
// provided string, if there is one.
func PublishStrategyFromString(s string) (PublishStrategy, error) {
	switch strings.ToLower(s) {
	case "commit":
		return PublishCommit, nil
	case "pull request":
		return PublishPullRequest, nil
	default:
		return "", fmt.Errorf("unsupported publish strategy: %s", s)
	}
}
