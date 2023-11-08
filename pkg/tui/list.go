package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ----------------------------------------------------------------------------
// styles for list
// ----------------------------------------------------------------------------

var (
	appStyle          = lipgloss.NewStyle().Padding(1, 2)
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle()
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	pickedItemStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	descrItemStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))

	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle   = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// model used by showList()
type model struct {
	list list.Model
	// text of value user picked/selected
	picked   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

// handling key-pressed events
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(list.Item)
			if ok {
				m.picked = i.FilterValue()
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// simple view function, for real rendering
// of items in list are responsible delegates
func (m model) View() string {
	return appStyle.Render(m.list.View())
}

// show list of items where items are rendered by given
// delegate
func showList(title string, items []list.Item, delegate list.ItemDelegate) (string, error) {

	l := list.New(items, delegate, 25, 20)
	l.Title = title
	l.SetShowHelp(true)
	l.SetShowStatusBar(false)
	l.SetShowPagination(true)

	p := tea.NewProgram(model{list: l}, tea.WithAltScreen())
	outModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if outModel.(model).quitting {
		return "", fmt.Errorf("terminating")
	}

	return outModel.(model).picked, nil
}
