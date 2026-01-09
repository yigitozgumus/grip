package selector

import (
	"os"

	"github.com/charmbracelet/lipgloss"
)

// Styles holds all the lipgloss styles for the selector.
type Styles struct {
	Container     lipgloss.Style
	Header        lipgloss.Style
	SearchBar     lipgloss.Style
	SearchPrompt  lipgloss.Style
	List          lipgloss.Style
	ItemNormal    lipgloss.Style
	ItemSelected  lipgloss.Style
	Cursor        lipgloss.Style
	FolderIcon    lipgloss.Style
	RepoName      lipgloss.Style
	RepoNameMatch lipgloss.Style
	Branch        lipgloss.Style
	Metadata      lipgloss.Style
	HelpBar       lipgloss.Style
	NoResults     lipgloss.Style
}

// DefaultStyles returns the default style configuration.
func DefaultStyles() Styles {
	// Use stderr for color detection (stdout may be piped by shell wrappers)
	renderer := lipgloss.NewRenderer(os.Stderr)

	subtle := lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"}
	highlight := lipgloss.AdaptiveColor{Light: "#7D56F4", Dark: "#AD8CFF"}
	cursor := lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	return Styles{
		Container: renderer.NewStyle().
			Padding(1, 2),

		Header: renderer.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			MarginBottom(1),

		SearchBar: renderer.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(0, 1).
			MarginBottom(1),

		SearchPrompt: renderer.NewStyle().
			Foreground(lipgloss.Color("205")),

		List: renderer.NewStyle().
			MarginBottom(1),

		ItemNormal: renderer.NewStyle().
			PaddingLeft(2),

		ItemSelected: renderer.NewStyle().
			Foreground(highlight),

		Cursor: renderer.NewStyle().
			Foreground(cursor).
			Bold(true),

		FolderIcon: renderer.NewStyle().
			Foreground(lipgloss.Color("220")),

		RepoName: renderer.NewStyle().
			Foreground(lipgloss.Color("252")),

		RepoNameMatch: renderer.NewStyle().
			Foreground(highlight).
			Bold(true),

		Branch: renderer.NewStyle().
			Foreground(lipgloss.Color("114")),

		Metadata: renderer.NewStyle().
			Foreground(subtle).
			Faint(true),

		HelpBar: renderer.NewStyle().
			Foreground(subtle).
			MarginTop(1),

		NoResults: renderer.NewStyle().
			Foreground(subtle).
			Italic(true).
			PaddingLeft(2),
	}
}
