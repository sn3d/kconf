package list

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ItemWithTitle interface {
	Title() string
}

type SimpleDelegate struct {
	picked int
	style  ItemStyle
}

func NewSimpleDelegate() *SimpleDelegate {
	return &SimpleDelegate{
		picked: 0,
		style:  DefaultStyle(),
	}
}

func (d *SimpleDelegate) Height() int { return 1 }

func (d *SimpleDelegate) Spacing() int { return 0 }

func (d *SimpleDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d *SimpleDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	ti, ok := listItem.(ItemWithTitle)
	if !ok {
		return
	}

	title := fmt.Sprintf("%d. %s", index+1, ti.Title())

	if d.IsPicked(index) {
		title = d.style.Picked.Render(title + "*")
	}

	if index == m.Index() {
		title = d.style.Selected.Render("  | " + title)
	} else {
		title = "    " + title
	}

	fmt.Fprint(w, title)
}

func (d *SimpleDelegate) Pick(index int) {
	d.picked = index
}

func (d *SimpleDelegate) IsPicked(index int) bool {
	return index == d.picked
}
