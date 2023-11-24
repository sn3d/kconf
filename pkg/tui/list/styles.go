package list

import "github.com/charmbracelet/lipgloss"

type ItemStyle struct {
	Picked      lipgloss.Style
	Selected    lipgloss.Style
	Description lipgloss.Style
}

func DefaultItemStyle() ItemStyle {
	return ItemStyle{
		Picked:      lipgloss.NewStyle().Foreground(lipgloss.Color("220")),
		Selected:    lipgloss.NewStyle().Foreground(lipgloss.Color("170")),
		Description: lipgloss.NewStyle().Foreground(lipgloss.Color("242")),
	}
}
