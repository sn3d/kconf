package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sn3d/kconf/pkg/tui/components/confirmation"
)

type Model struct {
	confirmation confirmation.Model
}

func NewModel() *Model {
	m := &Model{
		confirmation: confirmation.New("1", "Do you want to delete it?"),
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case confirmation.SubmittedMsg:
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.confirmation, cmd = m.confirmation.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.confirmation.View()
}

func main() {
	tea.NewProgram(NewModel(), tea.WithAltScreen()).Run()
}
