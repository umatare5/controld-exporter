// Package cli handles the execution of the CLI application.
package cli

import (
	"os/exec"
	"strings"
)

const (
	unknownVersion = "unknown"
)

// getVersion retrieves the Git tag or commit hash
func getVersion() string {
	if tag := getGitTag(); tag != "" {
		return tag
	}
	if hash := getGitCommitHash(); hash != "" {
		return hash
	}
	return unknownVersion
}

// getGitTag retrieves the latest Git tag
func getGitTag() string {
	return runCommand("git", "describe", "--tags", "--abbrev=0")
}

// getGitCommitHash retrieves the short Git commit hash
func getGitCommitHash() string {
	return runCommand("git", "rev-parse", "--short", "HEAD")
}

// runCommand executes a shell command and returns its output
func runCommand(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
