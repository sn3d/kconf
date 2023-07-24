package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sn3d/kconf/pkg/kconf"
)

type contextItem struct {
	Name       string
	ClusterURL string
	Namespace  string
}

func (i contextItem) FilterValue() string { return string(i.Name) }

type contextItemDelegate struct {
	picked string
}

func (d contextItemDelegate) Height() int                               { return 1 }
func (d contextItemDelegate) Spacing() int                              { return 0 }
func (d contextItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

// this is a main rendering procedure for each item
// in list
func (d contextItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(contextItem)
	if !ok {
		return
	}

	title := fmt.Sprintf("%d. %s", index+1, i.Name)

	var desc string
	if i.Namespace != "" {
		desc = fmt.Sprintf("url: %s ns: %s", i.ClusterURL, i.Namespace)
	} else {
		desc = fmt.Sprintf("url: %s", i.ClusterURL)
	}

	if i.Name == d.picked {
		title = pickedItemStyle.Render(title)
	} else {
		title = itemStyle.Render(title)
	}

	desc = descrItemStyle.Render(desc)

	if index == m.Index() {
		title = selectedItemStyle.Render("  | " + title)
		desc = selectedItemStyle.Render("  | " + desc)
	} else {
		title = "    " + title
		desc = "    " + desc
	}

	fmt.Fprint(w, title+"\n"+desc+"\n")
}

func ShowContextList(title string, picked string, conf *kconf.KubeConfig) (string, error) {
	items := []list.Item{}
	for _, ctx := range conf.Contexts {
		itm := contextItem{}
		itm.Name = ctx.Name

		cluster := conf.GetCluster(ctx.Context.Cluster)
		if cluster != nil {
			itm.ClusterURL = cluster.Cluster.Server
		}

		itm.Namespace = ctx.Context.Namespace
		items = append(items, itm)
	}

	return showList(title, items, contextItemDelegate{picked: picked})
}
