package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type stringItem string

func (i stringItem) FilterValue() string { return string(i) }

type stringItemDelegate struct {
	picked string
}

func (d stringItemDelegate) Height() int                               { return 1 }
func (d stringItemDelegate) Spacing() int                              { return 0 }
func (d stringItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d stringItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(stringItem)
	if !ok {
		return
	}

	title := fmt.Sprintf("%d. %s", index+1, i)
	if i == stringItem(d.picked) {
		title = pickedItemStyle.Render(title + "*")
	} else {
		title = itemStyle.Render(title)
	}

	if index == m.Index() {
		title = selectedItemStyle.Render("  | " + title)
	} else {
		title = "    " + title
	}

	fmt.Fprint(w, title)
}

func ShowSimpleList(title string, pickedID string, options []string) (string, error) {
	listItems := make([]list.Item, len(options))
	for i, opt := range options {
		listItems[i] = stringItem(opt)
	}

	return showList(title, listItems, stringItemDelegate{picked: pickedID})
}
