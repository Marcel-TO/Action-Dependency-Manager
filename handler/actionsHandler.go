// Package handler provides utilities for inspecting local GitHub Action
// repositories, such as discovering their workflow definitions.
package handler

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"marcel-to/actions-manager/logger"
)

// workflowsDir is the standard relative path where GitHub Actions workflow
// definitions live within a repository.
const workflowsDir = ".github/workflows"

// workflowExtensions holds the file extensions GitHub recognises as workflow
// files.
var workflowExtensions = map[string]struct{}{
	".yml":  {},
	".yaml": {},
}

// GetWorkflowFilePaths returns the paths of all GitHub Actions workflow files
// located in the .github/workflows directory of the repository rooted at
// repoPath.
//
// Only .yml and .yaml files directly inside .github/workflows are returned,
// matching GitHub's behaviour of ignoring nested directories.
// If the workflows directory does not exist, an empty slice and a nil error are returned.
func GetWorkflowFilePaths(logger *logger.Logger, repoPath string) ([]string, error) {
	dir := filepath.Join(repoPath, workflowsDir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read workflows directory %s: %w", dir, err)
	}

	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !isWorkflowFile(entry.Name()) {
			continue
		}
		paths = append(paths, filepath.Join(dir, entry.Name()))
	}

	return paths, nil
}

// isWorkflowFile reports whether name has a GitHub Actions workflow extension.
func isWorkflowFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	_, ok := workflowExtensions[ext]
	return ok
}
