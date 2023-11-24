package clst

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

	// convert clusters to items
	items := make([]bubblelist.Item, len(kc.AuthInfos))
	pickedIndex := -1
	for i := range kc.Clusters {
		items[i] = ClusterItem{
			Cluster: &kc.Clusters[i],
		}

		if kc.Clusters[i].Name == kc.GetCurrentContext().Context.Cluster {
			pickedIndex = i
		}
	}

	model := &Model{}
	model.kconfig = kc
	model.context = context

	model.list = list.New(items, list.NewSimpleDelegate())
	model.list.SetTitle(fmt.Sprintf("cluster for %s", context))
	model.list.Pick(pickedIndex)

	model.list.SetKeys(list.KeyMap{
		Pick: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "set cluster"),
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
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit without saving"),
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
		m.kconfig.ChangeCluster(m.context, msg.Picked)
		m.kconfig.Save()
		m.ExitMsg = msg
		return m, tea.Quit
	case list.RenameMsg:
		userItem := msg.Selected.(ClusterItem)
		m.kconfig.RenameCluster(userItem.Cluster.Name, msg.NewValue)
		return m, nil
	case list.SaveAndQuitMsg:
		m.kconfig.Save()
		m.ExitMsg = msg
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
