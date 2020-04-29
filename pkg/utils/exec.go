package utils

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// BinExists is a helper function which checks if a given binary
// exists somewhere on the PATH.
func BinExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Normalize is a helper function which strips a given []byte
// of any return characters (\r\n).
func Normalize(output string) string {
	s := strings.ReplaceAll(output, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

// RunCommand is a helper to run a command and collect the output from
// stdout and stderr.
func RunCommand(name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}
