package handler

import (
	"fmt"
	"os"
	"strings"

	"marcel-to/actions-manager/config"

	"gopkg.in/yaml.v2"
)

// usesKey is the workflow YAML key that references a GitHub Action dependency.
const usesKey = "uses"

// ParseActionDependencies reads the GitHub Actions workflow file at filePath and
// returns the action dependencies referenced by its `uses:` declarations.
//
// Each dependency is split into its action reference (e.g. "actions/checkout")
// and its version, which may be a tag, branch or commit hash (e.g. "7.0.1").
func ParseActionDependencies(filePath string) ([]config.ActionDependency, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file %s: %w", filePath, err)
	}

	var root interface{}
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("failed to parse workflow file %s: %w", filePath, err)
	}

	var refs []string
	collectUses(root, &refs)

	deps := make([]config.ActionDependency, 0, len(refs))
	seen := make(map[string]struct{}, len(refs))
	for _, ref := range refs {
		dep, ok := parseActionReference(ref)
		if !ok {
			continue
		}

		key := dep.Dependency + "@" + dep.Version
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		deps = append(deps, dep)
	}

	return deps, nil
}

// collectUses recursively walks a decoded YAML node and appends every string
// value associated with a `uses` key to refs.
func collectUses(node interface{}, refs *[]string) {
	switch n := node.(type) {
	case map[interface{}]interface{}:
		for key, value := range n {
			if name, ok := key.(string); ok && name == usesKey {
				if ref, ok := value.(string); ok {
					*refs = append(*refs, ref)
				}
			}
			collectUses(value, refs)
		}
	case map[string]interface{}:
		for key, value := range n {
			if key == usesKey {
				if ref, ok := value.(string); ok {
					*refs = append(*refs, ref)
				}
			}
			collectUses(value, refs)
		}
	case []interface{}:
		for _, item := range n {
			collectUses(item, refs)
		}
	}
}

// parseActionReference splits a `uses:` reference of the form
// "owner/repo@version" into its dependency and version parts. It reports false
// for references that are not versioned GitHub Action dependencies, such as
// local (./) or docker:// references.
func parseActionReference(ref string) (config.ActionDependency, bool) {
	ref = strings.TrimSpace(ref)

	// Drop any trailing inline comment or stray text. A valid action reference
	// never contains whitespace, so everything from the first space onward
	// (e.g. "@7.0.1 # Newest version") is not part of the dependency.
	if idx := strings.IndexAny(ref, " \t"); idx >= 0 {
		ref = ref[:idx]
	}

	if ref == "" ||
		strings.HasPrefix(ref, "./") ||
		strings.HasPrefix(ref, "../") ||
		strings.HasPrefix(ref, "docker://") {
		return config.ActionDependency{}, false
	}

	at := strings.LastIndex(ref, "@")
	if at <= 0 || at == len(ref)-1 {
		return config.ActionDependency{}, false
	}

	return config.ActionDependency{
		Dependency: ref[:at],
		Version:    ref[at+1:],
	}, true
}
