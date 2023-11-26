package ns

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	bubblelist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sn3d/kconf/pkg/kconf"
	"github.com/sn3d/kconf/pkg/tui/list"
)

type Model struct {
	list    list.Model
	kconfig *kconf.KubeConfig
	context string
	ExitMsg tea.Msg
}

func NewModel(kc *kconf.KubeConfig, context string) (*Model, error) {

	if context == "" {
		context = kc.CurrentContext
	}

	// convert namespaces to items
	namespaces, err := kc.GetAllNamespaces(context)
	if err != nil {
		return nil, err
	}

	items := make([]bubblelist.Item, len(namespaces))
	pickedIndex := -1
	for i := range namespaces {
		items[i] = NamespaceItem{
			Name: namespaces[i],
		}

		if namespaces[i] == kc.GetCurrentContext().Context.Namespace {
			pickedIndex = i
		}
	}

	model := &Model{}
	model.kconfig = kc
	model.context = context

	model.list = list.New(items, list.NewSimpleDelegate())
	model.list.SetTitle(fmt.Sprintf("namespace for %s", context))
	model.list.Pick(pickedIndex)

	model.list.SetKeys(list.KeyMap{
		Pick: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "set namespace"),
		),
		Rename: key.NewBinding(
			key.WithDisabled(),
		),
		SaveAndClose: key.NewBinding(
			key.WithDisabled(),
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
	return m.list.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case list.PickedMsg:
		return m, m.onPicked(msg)
	case list.CloseMsg:
		return m, m.onClose(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) onPicked(msg list.PickedMsg) tea.Cmd {
	m.kconfig.ChangeNamespace(m.context, msg.Picked)
	m.kconfig.Save()
	m.ExitMsg = msg
	return tea.Quit
}

func (m *Model) onClose(msg list.CloseMsg) tea.Cmd {
	m.ExitMsg = msg
	return tea.Quit
}
