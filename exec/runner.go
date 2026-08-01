package exec

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// EnsureGHInstalled verifies that the GitHub CLI ("gh") is available on the
// system PATH.
func EnsureGHInstalled() error {
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("the GitHub CLI ('gh') is required but was not found on your PATH; " +
			"install it from https://cli.github.com/ and run 'gh auth login'")
	}
	return nil
}

// Output runs the given command in dir and returns its trimmed stdout. An empty
// dir runs the command in the current working directory.
func Output(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command '%s %s' failed: %w", name, strings.Join(args, " "), wrapExitError(err))
	}
	return strings.TrimSpace(string(out)), nil
}

// Run executes the given command in the specified directory, streaming its
// output to stdout and stderr. Used for interactive commands whose output the
// user should be able to see
func Run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command '%s %s' failed: %w", name, strings.Join(args, " "), err)
	}
	return nil
}

// FetchLatestReleaseTag returns the tag name of the latest release of the given
// GitHub Action dependency (e.g. "Azure/functions-action").
func FetchLatestReleaseTag(dependency string) (string, error) {
	out, err := Output("", "gh", "release", "view", "--repo", dependency, "--json", "tagName", "--jq", ".tagName")
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest release tag for dependency %s: %w", dependency, err)
	}
	return out, nil
}

// FetchCommitHashFromTag returns the commit SHA that the given tag of the
// dependency resolves to.
func FetchCommitHashFromTag(dependency, tag string) (string, error) {
	out, err := Output("", "gh", "api", fmt.Sprintf("repos/%s/commits/%s", dependency, tag), "-H", "Accept: application/vnd.github.sha")
	if err != nil {
		return "", fmt.Errorf("failed to fetch commit hash for dependency %s and tag %s: %w", dependency, tag, err)
	}
	return out, nil
}

// wrapExitError enriches an *exec.ExitError with the stderr output that the
// failed command wrote, which exec.Cmd.Output otherwise leaves inaccessible.
func wrapExitError(err error) error {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && len(exitErr.Stderr) > 0 {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(exitErr.Stderr)))
	}
	return err
}
