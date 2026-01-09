package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	fuzzyfinder "github.com/ktr0731/go-fuzzyfinder"
	"github.com/yigitozgumus/grip/internal/config"
)

// ProjectItem represents a selectable project
type ProjectItem struct {
	Path string
	Repo *config.Repository
}

// SelectProject opens a fuzzy finder to select a project
func SelectProject(repos map[string]*config.Repository) (*ProjectItem, error) {
	if len(repos) == 0 {
		return nil, fmt.Errorf("no repositories found")
	}

	items := make([]ProjectItem, 0, len(repos))
	for path, repo := range repos {
		items = append(items, ProjectItem{
			Path: path,
			Repo: repo,
		})
	}

	idx, err := fuzzyfinder.Find(
		items,
		func(i int) string {
			name := filepath.Base(items[i].Path)
			branch := items[i].Repo.CurrentBranch
			return fmt.Sprintf("  📁 %-25s  ⎇ %s", name, branch)
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, _ int) string {
			if i == -1 {
				return ""
			}
			item := items[i]

			// Build a nicely formatted preview
			var sb strings.Builder
			sb.WriteString("╭─────────────────────────────────────╮\n")
			sb.WriteString("│  📂 PROJECT DETAILS                 │\n")
			sb.WriteString("╰─────────────────────────────────────╯\n\n")

			sb.WriteString(fmt.Sprintf("  📁 Name:     %s\n", filepath.Base(item.Path)))
			sb.WriteString(fmt.Sprintf("  🌿 Branch:   %s\n", item.Repo.CurrentBranch))
			sb.WriteString(fmt.Sprintf("  📍 Path:     %s\n\n", item.Path))

			sb.WriteString("╭─────────────────────────────────────╮\n")
			sb.WriteString("│  🌳 BRANCHES                        │\n")
			sb.WriteString("╰─────────────────────────────────────╯\n\n")

			sb.WriteString(fmt.Sprintf("  Local:  %d branches\n", len(item.Repo.LocalBranches)))
			sb.WriteString(fmt.Sprintf("  Remote: %d branches\n", len(item.Repo.RemoteBranches)))

			// Show first few local branches
			if len(item.Repo.LocalBranches) > 0 {
				sb.WriteString("\n  Local branches:\n")
				for j, b := range item.Repo.LocalBranches {
					if j >= 5 {
						sb.WriteString(fmt.Sprintf("    ... and %d more\n", len(item.Repo.LocalBranches)-5))
						break
					}
					marker := "  "
					if b == item.Repo.CurrentBranch {
						marker = "→ "
					}
					sb.WriteString(fmt.Sprintf("    %s%s\n", marker, b))
				}
			}

			return sb.String()
		}),
		fuzzyfinder.WithHeader("  ↑/↓: navigate  │  enter: select  │  esc: cancel  │  type to filter"),
		fuzzyfinder.WithPromptString("🔍 "),
	)

	if err != nil {
		return nil, err
	}

	return &items[idx], nil
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
