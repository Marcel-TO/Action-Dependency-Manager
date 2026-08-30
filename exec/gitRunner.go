package exec

import (
	"fmt"
	"strings"

	"marcel-to/actions-manager/config"
	"marcel-to/actions-manager/logger"
)

// maxListedDependencies caps how many changed dependencies are listed in the
// commit body; any beyond are summarised as "... (X more)"
const maxListedDependencies = 5

// CommitToGitHandler stages the workflow changes, commits and pushes them only if there are any changes to commit.
func CommitToGitHandler(log *logger.Logger, repoPath string, changed []config.ChangedDependency) error {
	if err := Run(repoPath, "git", "add", ".github/workflows"); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// `git status --porcelain` prints one line per change and nothing at all
	// when the tree is clean, so an empty result means there is nothing to commit.
	status, err := Output(repoPath, "git", "status", ".github/workflows", "--porcelain")
	if err != nil {
		return fmt.Errorf("failed to check repository status: %w", err)
	}
	if status == "" {
		log.Info("No changes to commit; repository already up to date.")
		return nil
	}

	if err := Run(repoPath, "git", "commit", "-m", buildCommitMessage(changed)); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}
	if err := Run(repoPath, "git", "push"); err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}

	return nil
}

// buildCommitMessage produces commit messages.
func buildCommitMessage(changed []config.ChangedDependency) string {
	if len(changed) == 0 {
		return "chore(deps): bump GitHub Action dependencies to latest releases"
	}

	noun := "dependencies"
	if len(changed) == 1 {
		noun = "dependency"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "chore(deps): bump %d GitHub Action %s\n\n", len(changed), noun)

	listed := changed
	if len(listed) > maxListedDependencies {
		listed = listed[:maxListedDependencies]
	}
	for _, c := range listed {
		fmt.Fprintf(&b, "- %s: %s -> %s\n", c.Dependency, c.OldVersion, c.NewVersion)
	}
	if remaining := len(changed) - len(listed); remaining > 0 {
		fmt.Fprintf(&b, "- ... (%d more)\n", remaining)
	}

	return b.String()
}

// PullFromRemote pulls the latest changes from the remote repository.
func PullFromRemote(repoPath string) error {
	err := Run(repoPath, "git", "pull")
	if err != nil {
		return fmt.Errorf("failed to pull latest changes from remote: %w", err)
	}

	return nil
}
