package confirmation

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	SelectYes key.Binding
	SelectNo  key.Binding
	Submit    key.Binding
	Abort     key.Binding
}

func DefaultKeyMap() *KeyMap {
	return &KeyMap{
		SelectYes: key.NewBinding(key.WithKeys("y", "left")),
		SelectNo:  key.NewBinding(key.WithKeys("n", "right")),
		Submit:    key.NewBinding(key.WithKeys("enter")),
		Abort:     key.NewBinding(key.WithKeys("ctrl+c", "esc")),
	}
}
