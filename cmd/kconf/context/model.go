package context

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	bubblelist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/sn3d/kconf/cmd/kconf/ns"
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/sn3d/kconf/pkg/tui/components/confirmation"
	"github.com/sn3d/kconf/pkg/tui/list"
)

type State int

const (
	ContextView State = iota
	NamespaceView
)

type Model struct {
	list    list.Model
	secList list.Model
	state   State
	kconfig *kconf.KubeConfig
	context string
	ExitMsg tea.Msg
}

func NewModel(kc *kconf.KubeConfig) (*Model, error) {

	model := &Model{}
	model.kconfig = kc

	// convert clusters to items
	items := make([]bubblelist.Item, len(kc.Contexts))
	pickedIndex := -1
	for i := range kc.Contexts {
		items[i] = ContextItem{
			Context: &kc.Contexts[i],
			Kconf:   kc,
		}

		if kc.Contexts[i].Name == kc.GetCurrentContext().Name {
			pickedIndex = i
		}
	}

	model.list = list.New(items, list.NewDescrDelegate())
	model.list.SetTitle(fmt.Sprintf("all contexts"))
	model.list.Pick(pickedIndex)

	model.list.SetKeys(list.KeyMap{
		Pick: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "set default"),
		),
		ChangeNs: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "set default namespace"),
		),
		Rename: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		SaveAndClose: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "save and quit"),
		),
		Close: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit"),
		),
	})

	return model, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.state == NamespaceView {
		return m.secList.View()
	}
	return m.list.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case ContextView:
		return m.updateContext(msg)
	case NamespaceView:
		return m.updateNamespace(msg)
	default:
		return m, nil
	}
}

func (m Model) updateContext(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case list.PickedMsg:
		return m, m.onPickedContext(msg)
	case list.RenameMsg:
		return m, m.onRenamed(msg)
	case confirmation.SubmittedMsg:
		return m, m.onDeleteContext(msg)
	case list.ChangeNsMsg:
		return m, m.onNamespaceView(msg)
	case list.CloseMsg:
		return m, m.onClose(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd

}

func (m Model) updateNamespace(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case list.PickedMsg:
		return m, m.onPickedNamespace(msg)
	case list.CloseMsg:
		return m, m.onCloseNamespace(msg)
	}

	var cmd tea.Cmd
	m.secList, cmd = m.secList.Update(msg)
	return m, cmd
}

func (m *Model) onNamespaceView(msg list.ChangeNsMsg) tea.Cmd {
	kc := m.kconfig

	// convert namespaces to items
	selectedItem := m.list.SelectedItem().(ContextItem)
	namespaces, err := kc.GetAllNamespaces(selectedItem.Context.Name)
	if err != nil {
		return nil
	}

	items := make([]bubblelist.Item, len(namespaces))
	pickedIndex := -1
	for i := range namespaces {
		items[i] = ns.NamespaceItem{
			Name: namespaces[i],
		}

		if namespaces[i] == selectedItem.Context.Context.Namespace {
			pickedIndex = i
		}
	}

	m.secList = list.New(items, list.NewSimpleDelegate())
	m.secList.SetTitle(fmt.Sprintf("namespaces for %s", selectedItem.Context.Name))
	m.secList.Pick(pickedIndex)

	m.secList.SetKeys(list.KeyMap{
		Pick:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select namespace")),
		Close: key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("q/esc", "back")),
	})

	m.state = NamespaceView

	return tea.ClearScreen
}

func (m *Model) SwitchToCtxView() {
	m.state = ContextView
}

func (m *Model) onDeleteContext(msg confirmation.SubmittedMsg) tea.Cmd {
	if msg.Confirmed {
		item := m.list.SelectedItem().(list.ItemWithTitle)
		m.kconfig.Remove(item.Title())
		m.list.RemoveSelected()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return cmd
}

func (m *Model) onPickedContext(msg list.PickedMsg) tea.Cmd {
	m.kconfig.CurrentContext = msg.Picked
	m.kconfig.Save()
	m.ExitMsg = msg
	return tea.Quit
}

func (m *Model) onRenamed(msg list.RenameMsg) tea.Cmd {
	item := msg.Selected.(ContextItem)
	m.kconfig.Rename(item.Context.Name, msg.NewValue)
	return nil
}

func (m *Model) onClose(msg list.CloseMsg) tea.Cmd {
	if !msg.Aborted {
		m.kconfig.Save()
	}
	m.ExitMsg = msg
	return tea.Quit
}

func (m *Model) onPickedNamespace(msg list.PickedMsg) tea.Cmd {
	ctxItem := m.list.SelectedItem().(ContextItem)
	m.kconfig.ChangeNamespace(ctxItem.Title(), msg.Picked)
	m.SwitchToCtxView()
	return tea.ClearScreen

}

func (m *Model) onCloseNamespace(msg list.CloseMsg) tea.Cmd {
	m.SwitchToCtxView()
	return tea.ClearScreen
}
