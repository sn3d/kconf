package list

import bubblelist "github.com/charmbracelet/bubbles/list"

type Picker interface {
	Pick(index int)
	IsPicked(index int) bool
}

type PickerDelegate interface {
	bubblelist.ItemDelegate
	Picker
}
