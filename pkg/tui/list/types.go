package list

import bubblelist "github.com/charmbracelet/bubbles/list"

type ItemWithTitle interface {
	Title() string
}

type ItemWithTitleAndDesr interface {
	ItemWithTitle
	Description() string
}

type Picker interface {
	Pick(index int)
	IsPicked(index int) bool
}

type PickerDelegate interface {
	bubblelist.ItemDelegate
	Picker
}
