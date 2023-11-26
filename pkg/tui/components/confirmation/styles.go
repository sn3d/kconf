package confirmation

import "github.com/charmbracelet/lipgloss"

type Style struct {
	Title    lipgloss.Style
	Active   lipgloss.Style
	Inactive lipgloss.Style
}

func DefaultStyle() *Style {
	return &Style{
		Title:    lipgloss.NewStyle().PaddingLeft(2),
		Active:   lipgloss.NewStyle().Foreground(lipgloss.Color("220")),
		Inactive: lipgloss.NewStyle().Foreground(lipgloss.Color("242")),
	}
}
