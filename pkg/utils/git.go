package utils

import "strings"

// IsDirty checks whether a repository is in a dirty state (has uncommitted changes).
func IsDirty() (bool, string) {
	out, err := RunCommand("git", "status", "--porcelain")
	if strings.TrimSpace(out) != "" || err != nil {
		return true, out
	}
	return false, ""
}

// GetTag gets the latest git tag for a repository.
func GetTag() (string, error) {
	out, err := RunCommand("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// InRepo checks whether the directory which chart-release is run from is
// a git repository.
func InRepo() bool {
	out, err := RunCommand("git", "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// GetCommit gets the latest short commit hash.
func GetCommit() (string, error) {
	out, err := RunCommand("git", "show", "--format='%h'", "HEAD", "-q")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// GetGitUserName gets the name of the git user, as specified in the git config.
func GetGitUserName() (string, error) {
	out, err := RunCommand("git", "config", "user.name")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// GetGitUserEmail gets the email address of the git user, as specified in the
// git config.
func GetGitUserEmail() (string, error) {
	out, err := RunCommand("git", "config", "user.email")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
