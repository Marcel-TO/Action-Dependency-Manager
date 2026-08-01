package operations

import (
	"marcel-to/actions-manager/config"
	"marcel-to/actions-manager/exec"
	"marcel-to/actions-manager/logger"
)

// resolvedRef is the latest release reference of a dependency: the commit SHA it
// resolves to and the tag that names that release.
type resolvedRef struct {
	sha string
	tag string
}

// resolveLatestReleases resolves every unique dependency in repoDeps to the
// latest release's commit SHA and tag, querying each dependency only once.
// Dependencies without a resolvable release are omitted from the result and
// reported with a warning rather than failing the caller.
func resolveLatestReleases(log *logger.Logger, repoDeps config.RepoDependencies) map[string]resolvedRef {
	resolved := make(map[string]resolvedRef)

	for _, workflow := range repoDeps.Workflows {
		for _, dep := range workflow.Dependencies {
			if _, done := resolved[dep.Dependency]; done {
				continue
			}

			tag, err := exec.FetchLatestReleaseTag(dep.Dependency)
			if err != nil {
				log.Warning("Skipping dependency without a resolvable release: " + dep.Dependency + " - " + err.Error())
				continue
			}

			sha, err := exec.FetchCommitHashFromTag(dep.Dependency, tag)
			if err != nil {
				log.Warning("Skipping dependency, failed to resolve commit hash: " + dep.Dependency + " - " + err.Error())
				continue
			}

			resolved[dep.Dependency] = resolvedRef{sha: sha, tag: tag}
			log.Debug("Resolved " + dep.Dependency + " -> " + sha + " (" + tag + ")")
		}
	}

	return resolved
}

// shortRef abbreviates a full 40-character commit SHA to its short form for
// display, leaving tags and other references unchanged.
func shortRef(ref string) string {
	if len(ref) == 40 {
		return ref[:7]
	}
	return ref
}
