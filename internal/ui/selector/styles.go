package selector

import "github.com/charmbracelet/lipgloss"

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
	subtle := lipgloss.AdaptiveColor{Light: "#666666", Dark: "#999999"}
	highlight := lipgloss.AdaptiveColor{Light: "#7D56F4", Dark: "#AD8CFF"}
	cursor := lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	return Styles{
		Container: lipgloss.NewStyle().
			Padding(1, 2),

		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).
			MarginBottom(1),

		SearchBar: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(0, 1).
			MarginBottom(1),

		SearchPrompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")),

		List: lipgloss.NewStyle().
			MarginBottom(1),

		ItemNormal: lipgloss.NewStyle().
			PaddingLeft(2),

		ItemSelected: lipgloss.NewStyle().
			Foreground(highlight),

		Cursor: lipgloss.NewStyle().
			Foreground(cursor).
			Bold(true),

		FolderIcon: lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")),

		RepoName: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),

		RepoNameMatch: lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true),

		Branch: lipgloss.NewStyle().
			Foreground(lipgloss.Color("114")),

		Metadata: lipgloss.NewStyle().
			Foreground(subtle).
			Faint(true),

		HelpBar: lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1),

		NoResults: lipgloss.NewStyle().
			Foreground(subtle).
			Italic(true).
			PaddingLeft(2),
	}
}
