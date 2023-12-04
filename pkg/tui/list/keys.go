package list

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Rename       key.Binding
	Pick         key.Binding
	ChangeNs     key.Binding
	Delete       key.Binding
	SaveAndClose key.Binding
	Close        key.Binding
	Help         key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Pick, k.Delete, k.SaveAndClose, k.Close, k.Help}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Pick, k.Rename, k.Delete},
		{k.ChangeNs},
		{k.SaveAndClose, k.Close},
	}
}
