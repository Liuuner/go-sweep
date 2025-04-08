package common

import "github.com/charmbracelet/bubbles/key"

type HotkeysType struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Help   key.Binding
	Quit   key.Binding
}

var Hotkeys = HotkeysType{
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h", "move left"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j", "move down"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k", "move up"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l", "move right"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "select"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("Q", "esc", "ctrl+c"),
		key.WithHelp("Q", "quit"),
	),
}
