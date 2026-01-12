package selector

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/yigitozgumus/grip/internal/ui/activity"
	"github.com/yigitozgumus/grip/internal/ui/common"
)

// View renders the selector UI.
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Title/Header
	b.WriteString(m.styles.Header.Render(m.title))
	b.WriteString("\n")

	// Search bar - simpler layout
	searchLine := m.styles.SearchPrompt.Render("Search: ") + m.searchInput.View()
	b.WriteString(searchLine)
	b.WriteString("\n\n")

	// List of repositories
	listContent := m.renderList()
	b.WriteString(listContent)

	// Help bar
	if m.showHelp {
		b.WriteString("\n")
		helpView := m.help.View(m.keymap)
		b.WriteString(m.styles.HelpBar.Render(helpView))
	}

	return b.String()
}

func (m Model) renderList() string {
	if len(m.filteredRepos) == 0 {
		return m.styles.NoResults.Render("No repositories found")
	}

	var lines []string

	// Calculate visible items based on height
	visibleCount := m.calculateVisibleItems()
	start, end := m.calculateViewport(visibleCount)

	for i := start; i < end; i++ {
		item := m.filteredRepos[i]
		isSelected := i == m.cursor
		lines = append(lines, m.renderItem(item, isSelected))
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderItem(item RepoItem, isSelected bool) string {
	// Cursor/selection indicator
	cursor := "  "
	if isSelected {
		cursor = m.styles.Cursor.Render("→ ")
	}

	// Folder icon
	icon := m.styles.FolderIcon.Render("📁")

	// Pad the raw name first, then apply styling
	const nameWidth = 30
	paddedName := item.Name
	if len(item.Name) < nameWidth {
		paddedName = item.Name + strings.Repeat(" ", nameWidth-len(item.Name))
	} else if len(item.Name) > nameWidth {
		paddedName = item.Name[:nameWidth-1] + "…"
	}

	// Apply highlighting to the padded name
	name := m.renderNameWithHighlight(paddedName, item.MatchedIndexes, isSelected)

	// Branch indicator (truncate if too long)
	branchName := item.CurrentBranch
	const maxBranchLen = 25
	if len(branchName) > maxBranchLen {
		branchName = branchName[:maxBranchLen-1] + "…"
	}
	branch := m.styles.Branch.Render(fmt.Sprintf("⎇ %s", branchName))

	// Right-aligned metadata: relative time + activity score
	relTime := common.FormatRelativeTime(item.LastUpdated)
	score := activity.FormatScore(item.ActivityScore)
	metadata := m.styles.Metadata.Render(fmt.Sprintf("%8s  %s", relTime, score))

	// Compose the line with proper spacing
	leftPart := cursor + icon + " " + name + "  " + branch

	// Calculate padding for right alignment
	totalWidth := m.width - 6 // Account for container padding
	leftWidth := lipgloss.Width(leftPart)
	metadataWidth := lipgloss.Width(metadata)
	padding := totalWidth - leftWidth - metadataWidth

	if padding < 2 {
		// Not enough space, just show left part
		return leftPart
	}

	return leftPart + strings.Repeat(" ", padding) + metadata
}

func (m Model) renderNameWithHighlight(name string, matchedIndexes []int, isSelected bool) string {
	if len(matchedIndexes) == 0 {
		if isSelected {
			return m.styles.ItemSelected.Render(name)
		}
		return m.styles.RepoName.Render(name)
	}

	// Build highlighted string
	var result strings.Builder
	matchSet := make(map[int]bool)
	for _, idx := range matchedIndexes {
		matchSet[idx] = true
	}

	for i, char := range name {
		s := string(char)
		if matchSet[i] {
			result.WriteString(m.styles.RepoNameMatch.Render(s))
		} else if isSelected {
			result.WriteString(m.styles.ItemSelected.Render(s))
		} else {
			result.WriteString(m.styles.RepoName.Render(s))
		}
	}

	return result.String()
}

func (m Model) calculateVisibleItems() int {
	// Reserve space for: header (2), search (3), help (2), padding
	reserved := 9
	available := m.height - reserved
	return max(available, 3)
}

func (m Model) calculateViewport(visibleCount int) (start, end int) {
	total := len(m.filteredRepos)

	if total <= visibleCount {
		return 0, total
	}

	// Keep cursor roughly centered
	halfVisible := visibleCount / 2
	start = max(m.cursor-halfVisible, 0)

	end = min(start+visibleCount, total)
	start = max(end-visibleCount, 0)

	return start, end
}
