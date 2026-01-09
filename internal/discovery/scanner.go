package discovery

import (
	"os"
	"path/filepath"

	"github.com/yigitozgumus/grip/internal/config"
	"github.com/yigitozgumus/grip/internal/git"
)

// ScanWorkspace discovers all git repositories in a workspace
func ScanWorkspace(ws config.Workspace) ([]*config.Repository, error) {
	var repos []*config.Repository

	err := scanDir(ws.Path, 0, ws.MaxDepth, func(path string) error {
		if git.IsGitRepo(path) {
			repo, err := git.RefreshRepoInfo(path)
			if err != nil {
				// Log error but continue scanning
				return nil
			}
			repos = append(repos, repo)
			return filepath.SkipDir // Don't descend into git repos
		}
		return nil
	})

	return repos, err
}

// ScanSingleRepo scans a single directory as a repository
func ScanSingleRepo(path string) (*config.Repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if !git.IsGitRepo(absPath) {
		return nil, os.ErrNotExist
	}

	return git.RefreshRepoInfo(absPath)
}

func scanDir(path string, currentDepth, maxDepth int, onDir func(string) error) error {
	if currentDepth > maxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil // Skip unreadable directories
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Skip hidden directories (but we need to check for .git inside dirs)
		if name[0] == '.' {
			continue
		}

		fullPath := filepath.Join(path, name)

		if err := onDir(fullPath); err != nil {
			if err == filepath.SkipDir {
				continue
			}
			return err
		}

		// Recurse into subdirectories
		if err := scanDir(fullPath, currentDepth+1, maxDepth, onDir); err != nil {
			return err
		}
	}

	return nil
}
