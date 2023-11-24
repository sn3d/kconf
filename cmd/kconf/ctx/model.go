package ctx

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	bubblelist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/sn3d/kconf/cmd/kconf/ns"
	"github.com/sn3d/kconf/pkg/kconf"
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
	items := make([]bubblelist.Item, len(kc.AuthInfos))
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
			key.WithHelp("enter", "set cluster"),
		),
		ChangeNs: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "change namespace"),
		),
		Rename: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename cluster"),
		),
		SaveAndQuit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "save and quit"),
		),
		Terminate: key.NewBinding(
			key.WithKeys("ctrl+c", "esc"),
			key.WithHelp("ctrl+c", "quit without saving"),
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
		m.kconfig.CurrentContext = msg.Picked
		m.kconfig.Save()
		m.ExitMsg = msg
		return m, tea.Quit
	case list.RenameMsg:
		item := msg.Selected.(ContextItem)
		m.kconfig.RenameContext(item.Context.Name, msg.NewValue)
		return m, nil
	case list.ChangeNsMsg:
		m.SwithToNsView()
		return m, tea.ClearScreen
	case list.SaveAndQuitMsg:
		m.kconfig.Save()
		m.ExitMsg = msg
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd

}

func (m Model) updateNamespace(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case list.PickedMsg:
		ctxItem := m.list.SelectedItem().(ContextItem)
		m.kconfig.ChangeNamespace(ctxItem.Title(), msg.Picked)
		m.SwitchToCtxView()
		return m, tea.ClearScreen
	case list.ExitMsg:
		m.SwitchToCtxView()
		return m, tea.ClearScreen
	}

	var cmd tea.Cmd
	m.secList, cmd = m.secList.Update(msg)
	return m, cmd
}

func (m *Model) SwithToNsView() {
	kc := m.kconfig

	// convert namespaces to items
	selectedItem := m.list.SelectedItem().(ContextItem)
	namespaces, err := kc.GetAllNamespaces(selectedItem.Context.Name)
	if err != nil {
		return
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
		Pick: key.NewBinding(key.WithKeys("enter")),
		Exit: key.NewBinding(key.WithKeys("esc", "q")),
	})

	m.state = NamespaceView
}

func (m *Model) SwitchToCtxView() {
	m.state = ContextView
}
