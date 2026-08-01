package handler

import (
	"fmt"

	"marcel-to/actions-manager/config"
	"marcel-to/actions-manager/logger"

	"gopkg.in/yaml.v2"
)

// HandleRepository discovers the GitHub Action dependencies used across every
// workflow file in the repository and returns them.
func HandleRepository(log *logger.Logger, repo config.RepoConfig) (config.RepoDependencies, error) {
	// Find all workflow files and extract their dependencies.
	actionPaths, err := GetWorkflowFilePaths(log, repo.RepoLocation)
	if err != nil {
		return config.RepoDependencies{}, err
	}

	workflows := make([]config.WorkflowDependencies, 0, len(actionPaths))
	for _, actionPath := range actionPaths {
		dependencies, err := ParseActionDependencies(actionPath)
		if err != nil {
			return config.RepoDependencies{}, err
		}
		workflows = append(workflows, config.WorkflowDependencies{
			FilePath:     actionPath,
			Dependencies: dependencies,
		})
	}

	repoDeps := config.RepoDependencies{
		Repo:      repo,
		Workflows: workflows,
	}

	// Print the discovered dependencies to stdout as YAML.
	data, err := yaml.Marshal(repoDeps)
	if err != nil {
		return config.RepoDependencies{}, fmt.Errorf("failed to marshal dependencies for repo %s: %w", repo.Name, err)
	}

	log.Debug(string(data))
	return repoDeps, nil
}
