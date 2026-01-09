package selector

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yigitozgumus/grip/internal/config"
)

// Run starts the selector UI and returns the result.
func Run(repos map[string]*config.Repository, opts Options) (Result, error) {
	if len(repos) == 0 {
		return Result{Cancelled: true}, nil
	}

	model := New(repos, opts)

	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return Result{Cancelled: true}, err
	}

	m := finalModel.(Model)

	return Result{
		Selected:  m.selected,
		Cancelled: m.cancelled,
	}, nil
}

// RunWithItems starts the selector with pre-converted items.
func RunWithItems(items []RepoItem, opts Options) (Result, error) {
	if len(items) == 0 {
		return Result{Cancelled: true}, nil
	}

	model := NewWithItems(items, opts)

	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return Result{Cancelled: true}, err
	}

	m := finalModel.(Model)

	return Result{
		Selected:  m.selected,
		Cancelled: m.cancelled,
	}, nil
}
