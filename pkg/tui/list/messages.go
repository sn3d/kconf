package list

import (
	bubblelist "github.com/charmbracelet/bubbles/list"
)

type PickedMsg struct {
	Picked string
}

type RenameMsg struct {
	Selected bubblelist.Item
	NewValue string
}

type CloseMsg struct {
	Aborted bool
}

type ChangeNsMsg struct {
	selected string
}
