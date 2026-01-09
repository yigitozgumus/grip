package selector

import (
	"path/filepath"
	"sort"
	"time"

	"github.com/yigitozgumus/grip/internal/config"
	"github.com/yigitozgumus/grip/internal/ui/activity"
)

// RepoItem represents a repository in the selector list.
type RepoItem struct {
	ID             string
	Name           string
	Path           string
	CurrentBranch  string
	LastUpdated    time.Time
	ActivityScore  float64
	MatchedIndexes []int
}

// convertRepos converts config repositories to RepoItems sorted by activity.
func convertRepos(repos map[string]*config.Repository) []RepoItem {
	items := make([]RepoItem, 0, len(repos))

	for path, repo := range repos {
		item := RepoItem{
			ID:            path,
			Name:          filepath.Base(path),
			Path:          path,
			CurrentBranch: repo.CurrentBranch,
			LastUpdated:   repo.LastUpdated,
			ActivityScore: activity.CalculateScore(repo),
		}
		items = append(items, item)
	}

	// Sort by activity score descending (most active first)
	sort.Slice(items, func(i, j int) bool {
		return items[i].ActivityScore > items[j].ActivityScore
	})

	return items
}

// repoSource implements fuzzy.Source for repo items.
type repoSource []RepoItem

func (r repoSource) String(i int) string {
	return r[i].Name
}

func (r repoSource) Len() int {
	return len(r)
}
