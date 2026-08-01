package operations

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"marcel-to/actions-manager/config"
	"marcel-to/actions-manager/exec"
	"marcel-to/actions-manager/handler"
	"marcel-to/actions-manager/logger"
)

// RunUpdate inspects the selected repositories and updates every Github Action dependency
// in their workflow files to the latest release, pinning them to the release's commit SHA
//
//	with the tag kept as a trailing comment. It optionally commits and pushes the changes.
func RunUpdate(log *logger.Logger, args config.UpdateArgumentConfig, cfg *config.Config) error {
	// Adjust debugging based on the provided argument.
	if args.Debug {
		log.SetDebugging(true)
	}

	log.Info("Starting update process...")
	selectedIds := make([]string, 0)

	if args.Selected != "" {
		selectedIds = strings.Split(args.Selected, ",")
	}
	for _, repo := range cfg.Repos {
		if args.Selected != "" && !slices.Contains(selectedIds, strconv.Itoa(repo.ID)) {
			log.Info("Skipping repository: " + repo.Name)
			continue
		}

		log.Info("")
		log.Info("")
		log.Info(strings.Repeat("=", 40))
		log.Info("Updating repository: " + repo.Name)
		log.Info(strings.Repeat("=", 40))

		// Pull the latest changes from remote if requested
		if args.Pull {
			log.Info("Pulling latest changes for repository...")
			err := exec.PullFromRemote(repo.RepoLocation)
			if err != nil {
				log.Error("Failed to pull latest changes for repository: " + repo.Name + " - " + err.Error())
				return err
			}
		}

		repoDeps, err := handler.HandleRepository(log, repo)
		if err != nil {
			log.Error("Failed to update repository: " + repo.Name + " - " + err.Error())
			return err
		}

		changed, err := updateDependencies(log, repoDeps)
		if err != nil {
			log.Error("Failed to update dependencies for repository: " + repo.Name + " - " + err.Error())
			return err
		}

		// Commit and push changes if requested
		if args.Commit {
			log.Info("Committing and pushing changes for repository...")
			err := exec.CommitToGitHandler(log, repo.RepoLocation, changed)
			if err != nil {
				log.Error("Failed to commit and push changes for repository: " + repo.Name + " - " + err.Error())
				return err
			}
		}
	}

	return nil
}

// updateDependencies pins every discovered GitHub Action dependency in the
// repository's workflow files to the latest release: the reference is rewritten
// to the release's commit SHA with the tag kept as a trailing comment
// (e.g. "uses: actions/checkout@<sha> # v4.2.2"). Dependencies without a
// resolvable release are skipped with a warning rather than failing the run.
func updateDependencies(log *logger.Logger, repoDeps config.RepoDependencies) ([]config.ChangedDependency, error) {
	// Resolve each unique dependency to its latest release tag and commit SHA.
	resolved := resolveLatestReleases(log, repoDeps)

	// Track which dependencies actually changed, deduplicated across all
	// workflow files: the same action is typically used in several files, so
	// keying by dependency keeps the resulting commit message concise.
	var changed []config.ChangedDependency
	seen := make(map[string]struct{})

	// Rewrite each workflow file, pinning every resolved dependency to its
	// latest commit SHA with the tag as a trailing comment.
	for _, workflow := range repoDeps.Workflows {
		content, err := os.ReadFile(workflow.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read workflow file %s: %w", workflow.FilePath, err)
		}

		updated := string(content)
		for _, dep := range workflow.Dependencies {
			ref, ok := resolved[dep.Dependency]
			if !ok {
				continue
			}

			// Record the dependency as changed (once) when its current pin
			// differs from the resolved latest release commit.
			if _, ok := seen[dep.Dependency]; !ok && dep.Version != ref.sha {
				seen[dep.Dependency] = struct{}{}
				changed = append(changed, config.ChangedDependency{
					Dependency: dep.Dependency,
					OldVersion: shortRef(dep.Version),
					NewVersion: ref.tag,
				})
			}

			// Match a `uses:` line referencing this dependency and replace the
			// whole reference (and any existing trailing comment) while keeping
			// the original indentation captured in group 1. The trailing part is
			// matched with [^\r\n]* (not .*) so a CRLF file keeps its \r: Go's
			// `.` matches \r, which would otherwise strip it and leave the line
			// LF-terminated, producing mixed line endings that git warns about.
			pattern := regexp.MustCompile(`(?m)^(\s*(?:-\s*)?uses:\s*)` + regexp.QuoteMeta(dep.Dependency) + `@\S+[^\r\n]*`)
			replacement := "${1}" + dep.Dependency + "@" + ref.sha + " # " + ref.tag
			updated = pattern.ReplaceAllString(updated, replacement)
		}

		if updated == string(content) {
			log.Debug("No changes for workflow file: " + workflow.FilePath)
			continue
		}

		if err := os.WriteFile(workflow.FilePath, []byte(updated), 0o644); err != nil {
			return nil, fmt.Errorf("failed to write workflow file %s: %w", workflow.FilePath, err)
		}
		log.Info("Updated workflow file: " + workflow.FilePath)
	}

	return changed, nil
}
