package list

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ItemWithTitleAndDesr interface {
	ItemWithTitle
	Description() string
}

type DescrDelegate struct {
	picked int
	style  ItemStyle
}

func NewDescrDelegate() *DescrDelegate {
	return &DescrDelegate{
		picked: -1,
		style:  DefaultStyle(),
	}
}

func (d *DescrDelegate) Height() int { return 1 }

func (d *DescrDelegate) Spacing() int { return 0 }

func (d *DescrDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d *DescrDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(ItemWithTitleAndDesr)
	if !ok {
		return
	}

	title := fmt.Sprintf("%d. %s", index+1, item.Title())
	descr := d.style.Description.Render(fmt.Sprintf("   %s", item.Description()))

	if d.IsPicked(index) {
		title = d.style.Picked.Render(title + "*")
	}

	if index == m.Index() {
		title = d.style.Selected.Render("  | " + title)
		descr = d.style.Selected.Render("  | " + descr)
	} else {
		title = "    " + title
		descr = "    " + descr
	}

	fmt.Fprintf(w, "%s\n%s\n", title, descr)
}

func (d *DescrDelegate) IsPicked(index int) bool {
	return d.picked == index
}

func (d *DescrDelegate) Pick(index int) {
	d.picked = index
}
