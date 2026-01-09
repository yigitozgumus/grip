package selector

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sahilm/fuzzy"
)

// Update handles messages and updates the model state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.Quit):
			m.quitting = true
			m.cancelled = true
			return m, tea.Quit

		case key.Matches(msg, m.keymap.Cancel):
			m.cancelled = true
			return m, tea.Quit

		case key.Matches(msg, m.keymap.Up):
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case key.Matches(msg, m.keymap.Down):
			if m.cursor < len(m.filteredRepos)-1 {
				m.cursor++
			}
			return m, nil

		case key.Matches(msg, m.keymap.Select):
			if len(m.filteredRepos) > 0 {
				selected := m.filteredRepos[m.cursor]
				m.selected = &selected
			}
			return m, tea.Quit

		case key.Matches(msg, m.keymap.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
	}

	// Update search input
	prevSearch := m.searchInput.Value()
	m.searchInput, cmd = m.searchInput.Update(msg)

	// If search changed, re-filter
	if m.searchInput.Value() != prevSearch {
		m.filterRepos()
		m.cursor = 0
	}

	return m, cmd
}

// filterRepos applies fuzzy filtering based on search input.
func (m *Model) filterRepos() {
	query := m.searchInput.Value()

	if query == "" {
		// No filter, show all repos sorted by activity
		m.filteredRepos = make([]RepoItem, len(m.repos))
		copy(m.filteredRepos, m.repos)
		// Clear match indexes
		for i := range m.filteredRepos {
			m.filteredRepos[i].MatchedIndexes = nil
		}
		return
	}

	// Fuzzy match against repo names
	matches := fuzzy.FindFrom(query, repoSource(m.repos))

	m.filteredRepos = make([]RepoItem, 0, len(matches))
	for _, match := range matches {
		item := m.repos[match.Index]
		item.MatchedIndexes = match.MatchedIndexes
		m.filteredRepos = append(m.filteredRepos, item)
	}
}
