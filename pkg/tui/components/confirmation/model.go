package confirmation

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	id        string
	title     string
	confirmed bool
	keys      *KeyMap
	styles    *Style
}

func New(id string, title string) Model {
	m := Model{}
	m.id = id
	m.title = title
	m.keys = DefaultKeyMap()
	m.styles = DefaultStyle()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.SelectYes):
			m.confirmed = true
			return m, nil
		case key.Matches(msg, m.keys.SelectNo):
			m.confirmed = false
			return m, nil
		case key.Matches(msg, m.keys.Submit):
			return m, Submit(m.id, m.confirmed)
		case key.Matches(msg, m.keys.Abort):
			return m, Submit(m.id, false)
		}
	}
	return m, nil
}

func (m Model) View() string {
	yesRender := "Yes"
	if m.confirmed {
		yesRender = m.styles.Active.Render(">" + yesRender)
	} else {
		yesRender = m.styles.Inactive.Render(" " + yesRender)
	}

	noRender := "No"
	if !m.confirmed {
		noRender = m.styles.Active.Render(">" + noRender)
	} else {
		noRender = m.styles.Inactive.Render(" " + noRender)
	}

	titleRender := m.styles.Title.Render(m.title)

	render := fmt.Sprintf("%s: %s %s\n", titleRender, yesRender, noRender)
	return render
}

func (m *Model) IsConfirmed() bool {
	return m.confirmed
}
