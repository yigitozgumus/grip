package selector

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yigitozgumus/grip/internal/config"
)

// Result represents the outcome of the selector.
type Result struct {
	Selected  *RepoItem
	Cancelled bool
}

// Options configures the selector behavior.
type Options struct {
	Title    string
	ShowHelp bool
}

// Model is the main Bubble Tea model for the repository selector.
type Model struct {
	// Configuration
	title string
	repos []RepoItem

	// State
	filteredRepos []RepoItem
	cursor        int
	selected      *RepoItem
	quitting      bool
	cancelled     bool

	// Components
	searchInput textinput.Model
	help        help.Model
	keymap      KeyMap
	styles      Styles

	// Dimensions
	width  int
	height int

	// Options
	showHelp bool
}

// New creates a new selector model from config repositories.
func New(repos map[string]*config.Repository, opts Options) Model {
	items := convertRepos(repos)
	return newWithItems(items, opts)
}

// NewWithItems creates a selector model with pre-converted items.
func NewWithItems(items []RepoItem, opts Options) Model {
	return newWithItems(items, opts)
}

func newWithItems(items []RepoItem, opts Options) Model {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Prompt = ""
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	h := help.New()
	h.ShowAll = false

	title := opts.Title
	if title == "" {
		title = "Select Repository"
	}

	return Model{
		title:         title,
		repos:         items,
		filteredRepos: items,
		cursor:        0,
		searchInput:   ti,
		help:          h,
		keymap:        DefaultKeyMap(),
		styles:        DefaultStyles(),
		showHelp:      opts.ShowHelp,
		width:         80,
		height:        24,
	}
}

// Init initializes the model.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Selected returns the currently selected item, if any.
func (m Model) Selected() *RepoItem {
	return m.selected
}

// Cancelled returns whether the selection was cancelled.
func (m Model) Cancelled() bool {
	return m.cancelled
}
