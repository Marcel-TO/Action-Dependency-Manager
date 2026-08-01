package operations

import (
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"marcel-to/actions-manager/config"
	"marcel-to/actions-manager/exec"
	"marcel-to/actions-manager/handler"
	"marcel-to/actions-manager/logger"
)

// RunCheck inspects the selected repositories and reports, per workflow file,
// which GitHub Action dependencies have a newer release available and which are
// already pinned to the latest one. It is read-only and never modifies any
// files.
func RunCheck(log *logger.Logger, args config.CheckArgumentConfig, cfg *config.Config) error {
	// Adjust debugging based on the provided argument.
	if args.Debug {
		log.SetDebugging(true)
	}

	log.Info("Starting check process...")
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
		log.Info("Checking repository: " + repo.Name)
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
			log.Error("Failed to inspect repository: " + repo.Name + " - " + err.Error())
			return err
		}

		checkRepository(log, repoDeps)
	}

	return nil
}

// checkRepository resolves the latest release of every dependency in the
// repository and logs, per workflow file, whether each dependency is up to date
// or has a newer version available, followed by a short summary.
func checkRepository(log *logger.Logger, repoDeps config.RepoDependencies) {
	resolved := resolveLatestReleases(log, repoDeps)

	var outdated, upToDate int
	for _, workflow := range repoDeps.Workflows {
		log.Info("")
		log.Info("  " + filepath.Base(workflow.FilePath))
		log.Info("  " + strings.Repeat("-", len(filepath.Base(workflow.FilePath))))

		for _, dep := range workflow.Dependencies {
			ref, ok := resolved[dep.Dependency]
			if !ok {
				// The dependency has no resolvable release
				continue
			}

			if dep.Version == ref.sha {
				upToDate++
				log.Info("    ✓ up to date        " + dep.Dependency + " (" + ref.tag + ")")
				continue
			}

			outdated++
			log.Info("    ↑ update available  " + dep.Dependency + "  " +
				shortRef(dep.Version) + " -> " + shortRef(ref.sha) + " (" + ref.tag + ")")
		}
	}

	log.Info("")

	log.Info(fmt.Sprintf("Check summary for %s: %d update(s) available, %d up to date.", repoDeps.Repo.Name, outdated, upToDate))
}
