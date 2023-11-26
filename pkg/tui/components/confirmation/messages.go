package confirmation

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SubmittedMsg struct {
	ID        string
	Confirmed bool
}

func Submit(id string, confirmed bool) tea.Cmd {
	return func() tea.Msg {
		return SubmittedMsg{
			ID:        id,
			Confirmed: confirmed,
		}
	}
}
