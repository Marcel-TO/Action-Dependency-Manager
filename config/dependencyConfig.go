package config

// ActionDependency is a single GitHub Action dependency and its pinned version,
// which may be a tag, branch or commit hash.
type ActionDependency struct {
	Dependency string `yaml:"dependency"`
	Version    string `yaml:"version"`
}

// WorkflowDependencies holds the GitHub Action dependencies discovered in a
// single workflow file.
type WorkflowDependencies struct {
	FilePath     string             `yaml:"filePath"`
	Dependencies []ActionDependency `yaml:"dependencies"`
}

// RepoDependencies is the set of GitHub Action dependencies discovered across
// all workflow files of a single repository.
type RepoDependencies struct {
	Repo      RepoConfig
	Workflows []WorkflowDependencies `yaml:"workflows"`
}

// ChangedDependency describes a single GitHub Action dependency that was
// re-pinned during an update, used to build an informative commit message.
type ChangedDependency struct {
	Dependency string
	OldVersion string
	NewVersion string
}
