package list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	bubblelist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sn3d/kconf/pkg/tui/components/confirmation"
)

type State int

const (
	ListingState State = iota
	RenameState
	DeleteState
)

type Model struct {
	state State

	list         bubblelist.Model
	prompt       textinput.Model
	confirmation confirmation.Model
	help         help.Model

	listPicker Picker

	Keys KeyMap
}

func New(items []bubblelist.Item, delegate PickerDelegate) Model {
	model := Model{}
	model.listPicker = delegate

	model.list = bubblelist.New(items, delegate, 40, 20)
	model.list.SetShowStatusBar(false)
	model.list.SetShowHelp(false)
	model.list.DisableQuitKeybindings()

	model.prompt = textinput.New()
	model.prompt.CharLimit = 124

	model.help = help.New()

	return model
}

func (m *Model) SelectedItem() bubblelist.Item {
	return m.list.SelectedItem()
}

func (m *Model) RemoveSelected() {
	index := m.list.Index()
	m.list.RemoveItem(index)
	m.list.Select(-1)
}

func (m *Model) SetTitle(title string) {
	m.list.Title = title
}

func (m *Model) SetKeys(keys KeyMap) {
	m.Keys = keys
	m.Keys.Help = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "more"),
	)
}

func (m *Model) Pick(index int) {
	m.listPicker.Pick(index)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	switch m.state {
	case RenameState:
		return m.list.View() + "\n" + m.prompt.View()
	case DeleteState:
		return m.list.View() + "\n" + m.confirmation.View()
	default:
		help := lipgloss.NewStyle().PaddingLeft(2).Render(m.help.View(m.Keys))
		return m.list.View() + "\n" + help
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.prompt.Width = msg.Width
		m.help.Width = msg.Width
		return m, nil
	}

	switch m.state {
	case RenameState:
		return m.updateRename(msg)
	case DeleteState:
		return m.updateDeletion(msg)
	default:
		return m.updateListing(msg)
	}
}

func (m Model) updateListing(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Rename):
			return m, m.onRename()
		case key.Matches(msg, m.Keys.Delete):
			return m, m.onDelete()
		case key.Matches(msg, m.Keys.Help):
			return m, m.onHelp()
		case key.Matches(msg, m.Keys.ChangeNs):
			return m, m.onChangeNs()
		case key.Matches(msg, m.Keys.Pick):
			return m, m.onPick()
		case key.Matches(msg, m.Keys.Close):
			return m, m.onClose()
		case key.Matches(msg, m.Keys.SaveAndClose):
			return m, m.onSaveAndClose()
		case msg.Type == tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) updateRename(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, m.onRenameConfirmed()
		case tea.KeyEsc:
			m.toListing()
			return m, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.prompt, cmd = m.prompt.Update(msg)
	return m, cmd
}

func (m Model) updateDeletion(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case confirmation.SubmittedMsg:
		m.toListing()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.confirmation, cmd = m.confirmation.Update(msg)
	return m, cmd

}

func (m *Model) toListing() tea.Cmd {
	m.state = ListingState
	return nil
}

func (m *Model) onRename() tea.Cmd {
	m.state = RenameState
	selectedItem := m.list.Items()[m.list.Index()].(ItemWithTitle)
	m.prompt.Placeholder = selectedItem.Title()
	m.prompt.SetValue("")
	m.prompt.Focus()

	m.list.SetHeight(20)
	m.help.ShowAll = false

	return nil
}

func (m *Model) onRenameConfirmed() tea.Cmd {
	m.toListing()
	value := m.prompt.Value()
	if value != "" {
		selectedItem := m.list.Items()[m.list.Index()]
		cmd := func() tea.Msg {
			return RenameMsg{
				Selected: selectedItem,
				NewValue: value,
			}
		}
		return cmd
	} else {
		return nil
	}
}

func (m *Model) onDelete() tea.Cmd {
	m.state = DeleteState
	selectedItem := m.list.Items()[m.list.Index()].(ItemWithTitle)
	msg := fmt.Sprintf("Do you want to delete '%s' ?", selectedItem.Title())
	m.confirmation = confirmation.New("DELETE", msg)

	m.list.SetHeight(20)
	m.help.ShowAll = false

	return nil
}

func (m *Model) onHelp() tea.Cmd {
	m.help.ShowAll = !m.help.ShowAll
	if m.help.ShowAll {
		m.list.SetHeight(m.list.Height() - 2)
	} else {
		m.list.SetHeight(m.list.Height() + 2)
	}
	return nil
}

func (m *Model) onPick() tea.Cmd {
	index := m.list.Index()
	if index >= 0 {
		m.listPicker.Pick(index)
		pickedItem := m.list.Items()[index]
		pickedTitle := pickedItem.(ItemWithTitle).Title()
		cmd := func() tea.Msg {
			return PickedMsg{
				Picked: pickedTitle,
			}
		}
		return cmd
	} else {
		return nil
	}
}

func (m *Model) onChangeNs() tea.Cmd {
	pickedItem := m.list.Items()[m.list.Index()]
	title := pickedItem.(ItemWithTitle).Title()
	cmd := func() tea.Msg {
		return ChangeNsMsg{
			selected: title,
		}
	}
	return cmd
}

func (m *Model) onClose() tea.Cmd {
	return func() tea.Msg {
		return CloseMsg{
			Aborted: true,
		}
	}
}

func (m *Model) onSaveAndClose() tea.Cmd {
	return func() tea.Msg {
		return CloseMsg{
			Aborted: false,
		}
	}
}
