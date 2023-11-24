package list

import (
	"github.com/charmbracelet/bubbles/key"
	bubblelist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	ListingState State = iota
	PromptState
)

type Model struct {
	state State

	list bubblelist.Model

	listPicker Picker
	prompt     textinput.Model

	Keys KeyMap
}

func New(items []bubblelist.Item, delegate PickerDelegate) Model {
	model := Model{}
	model.listPicker = delegate

	model.list = bubblelist.New(items, delegate, 40, 20)
	model.list.SetShowStatusBar(false)
	model.list.SetShowHelp(true)
	model.list.DisableQuitKeybindings()

	model.prompt = textinput.New()
	model.prompt.CharLimit = 124

	return model
}

func (m *Model) SetTitle(title string) {
	m.list.Title = title
}

func (m *Model) SetKeys(keys KeyMap) {
	m.Keys = keys
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.ChangeNs,
			keys.Pick,
			keys.Rename,
			keys.SaveAndQuit,
			keys.Terminate,
		}
	}
}

func (m *Model) Pick(index int) {
	m.listPicker.Pick(index)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	switch m.state {
	case PromptState:
		return m.list.View() + "\n" + m.prompt.View()
	default:
		return m.list.View()
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.prompt.Width = msg.Width
		return m, nil
	}

	switch m.state {
	case PromptState:
		return m.updatePrompt(msg)
	default:
		return m.updateListing(msg)
	}
}

func (m Model) updateListing(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Rename):
			selectedItem := m.list.Items()[m.list.Index()].(ItemWithTitle)
			m.ShowPrompt(selectedItem.Title())
			return m, nil
		case key.Matches(msg, m.Keys.ChangeNs):
			pickedItem := m.list.Items()[m.list.Index()]
			title := pickedItem.(ItemWithTitle).Title()
			return m, ChangeNsCmd(title)
		case key.Matches(msg, m.Keys.Pick):
			m.listPicker.Pick(m.list.Index())
			pickedItem := m.list.Items()[m.list.Index()]
			return m, PickCmd(pickedItem.(ItemWithTitle).Title())
		case key.Matches(msg, m.Keys.Terminate):
			return m, tea.Quit
		case key.Matches(msg, m.Keys.SaveAndQuit):
			return m, SaveAndQuit
		case key.Matches(msg, m.Keys.Exit):
			return m, ExitCmd
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) updatePrompt(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.closePrompt()
			value := m.prompt.Value()
			if value != "" {
				selectedItem := m.list.Items()[m.list.Index()]
				return m, RenameCmd(selectedItem, value)
			} else {
				return m, nil
			}
		case tea.KeyEsc:
			m.closePrompt()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.prompt, cmd = m.prompt.Update(msg)
	return m, cmd
}

func (m *Model) ShowPrompt(placeholder string) {
	m.state = PromptState
	m.list.SetShowHelp(false)
	m.list.SetHeight(m.list.Height() - 1)

	m.prompt.Placeholder = placeholder
	m.prompt.SetValue("")
	m.prompt.Focus()
}

func (m *Model) closePrompt() {
	m.state = ListingState
	m.list.SetShowHelp(true)
	m.list.SetShowStatusBar(false)
	m.list.SetHeight(m.list.Height() + 1)
}

func (m *Model) SelectedItem() bubblelist.Item {
	return m.list.SelectedItem()
}
