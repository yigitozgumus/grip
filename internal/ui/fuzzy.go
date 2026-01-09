package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/yigitozgumus/grip/internal/config"
	"github.com/yigitozgumus/grip/internal/ui/selector"
)

// ProjectItem represents a selectable project
type ProjectItem struct {
	Path string
	Repo *config.Repository
}

// SelectProject opens a Bubble Tea TUI to select a project
func SelectProject(repos map[string]*config.Repository) (*ProjectItem, error) {
	if len(repos) == 0 {
		return nil, fmt.Errorf("no repositories found")
	}

	result, err := selector.Run(repos, selector.Options{
		Title:    "Select Project",
		ShowHelp: true,
	})

	if err != nil {
		return nil, err
	}

	if result.Cancelled || result.Selected == nil {
		return nil, fmt.Errorf("selection cancelled")
	}

	// Find the original repo from the config
	repo, ok := repos[result.Selected.Path]
	if !ok {
		return nil, fmt.Errorf("repository not found: %s", result.Selected.Path)
	}

	return &ProjectItem{
		Path: result.Selected.Path,
		Repo: repo,
	}, nil
}

// BranchItem represents a selectable branch
type BranchItem struct {
	Name      string
	IsRemote  bool
	IsCurrent bool
}

// SelectBranch opens a fuzzy finder to select a branch
func SelectBranch(repo *config.Repository) (string, error) {
	items := make([]BranchItem, 0)

	// Add local branches first
	for _, branch := range repo.LocalBranches {
		items = append(items, BranchItem{
			Name:      branch,
			IsRemote:  false,
			IsCurrent: branch == repo.CurrentBranch,
		})
	}

	// Add remote branches
	for _, branch := range repo.RemoteBranches {
		// Skip HEAD references
		if strings.HasSuffix(branch, "/HEAD") {
			continue
		}
		items = append(items, BranchItem{
			Name:     branch,
			IsRemote: true,
		})
	}

	if len(items) == 0 {
		return "", fmt.Errorf("no branches found")
	}

	repoName := filepath.Base(repo.Path)

	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			item := items[i]

			if item.IsCurrent {
				return fmt.Sprintf("  ● %s  ← current", item.Name)
			}

			if item.IsRemote {
				return fmt.Sprintf("  ☁ %s", item.Name)
			}

			return fmt.Sprintf("  ○ %s", item.Name)
		},
		fuzzyfinder.WithHeader(fmt.Sprintf("  📁 %s  │  ● local  ☁ remote  │  type to filter", repoName)),
		fuzzyfinder.WithPromptString("🌿 "),
	)

	if err != nil {
		return "", err
	}

	selected := items[idx]

	// For remote branches, extract the branch name (origin/main -> main)
	branchName := selected.Name
	if selected.IsRemote && strings.Contains(branchName, "/") {
		parts := strings.SplitN(branchName, "/", 2)
		if len(parts) == 2 {
			branchName = parts[1]
		}
	}

	return branchName, nil
}
