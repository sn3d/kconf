package user

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

	// convert auth infos(users) to list items
	items := make([]bubblelist.Item, len(kc.AuthInfos))
	pickedItem := -1
	for i := range kc.AuthInfos {
		items[i] = UserItem{
			User: &kc.AuthInfos[i],
		}

		if kc.AuthInfos[i].Name == kc.GetCurrentContext().Context.AuthInfo {
			pickedItem = i
		}
	}

	model := &Model{}
	model.kconfig = kc
	model.context = context

	model.list = list.New(items, list.NewSimpleDelegate())
	model.list.SetTitle(fmt.Sprintf("user for %s", context))
	model.list.Pick(pickedItem)

	model.list.SetKeys(list.KeyMap{
		Pick: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "set user"),
		),
		Rename: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename user"),
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
	return m.list.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case list.PickedMsg:
		return m, m.onPicked(msg)
	case list.RenameMsg:
		return m, m.onRename(msg)
	case list.CloseMsg:
		return m, m.onClose(msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) onPicked(msg list.PickedMsg) tea.Cmd {
	m.kconfig.ChangeUser(m.context, msg.Picked)
	m.kconfig.Save()
	m.ExitMsg = msg
	return tea.Quit
}

func (m *Model) onRename(msg list.RenameMsg) tea.Cmd {
	userItem := msg.Selected.(UserItem)
	m.kconfig.RenameUser(userItem.User.Name, msg.NewValue)
	return nil
}

func (m *Model) onClose(msg list.CloseMsg) tea.Cmd {
	if !msg.Aborted {
		m.kconfig.Save()
	}
	m.ExitMsg = msg
	return tea.Quit
}
