package list

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Rename      key.Binding
	Pick        key.Binding
	ChangeNs    key.Binding
	SaveAndQuit key.Binding
	Exit        key.Binding
	Terminate   key.Binding
}
