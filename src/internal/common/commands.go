package common

import tea "github.com/charmbracelet/bubbletea"

func asdf() {
	// This is a placeholder function to prevent the package from being empty.
	// It can be removed or replaced with actual code later.
	tea.Quit()
}

// QuitGameMsg signals that the program to close the game. You can send a [QuitGameMsg] with
// [QuitGame].
type QuitGameMsg struct{}

// QuitGame is a special command that tells the program to close the current game.
func QuitGame() tea.Msg {
	return QuitGameMsg{}
}
