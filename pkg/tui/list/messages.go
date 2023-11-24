package list

import (
	bubblelist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type PickedMsg struct {
	Picked string
}

func PickCmd(picked string) tea.Cmd {
	return func() tea.Msg {
		return PickedMsg{
			Picked: picked,
		}
	}
}

type RenameMsg struct {
	Selected bubblelist.Item
	NewValue string
}

func RenameCmd(selected bubblelist.Item, newValue string) tea.Cmd {
	return func() tea.Msg {
		return RenameMsg{
			Selected: selected,
			NewValue: newValue,
		}
	}
}

type SaveAndQuitMsg struct{}

func SaveAndQuit() tea.Msg {
	return SaveAndQuitMsg{}
}

type ExitMsg struct{}

func ExitCmd() tea.Msg {
	return ExitMsg{}
}

type ChangeNsMsg struct {
	selected string
}

func ChangeNsCmd(selected string) tea.Cmd {
	return func() tea.Msg {
		return ChangeNsMsg{
			selected: selected,
		}
	}
}
