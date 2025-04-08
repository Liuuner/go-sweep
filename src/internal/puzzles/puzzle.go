package puzzles

import (
	"github.com/Liuuner/go-puzzles/src/internal/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Puzzle interface {
	New() Puzzle
	Init() tea.Cmd
	View() string
	Update(msg tea.Msg) (Puzzle, tea.Cmd)
	Preview() string
}

type EmptyPuzzle struct{}

func (p EmptyPuzzle) New() Puzzle {
	return p
}

func (EmptyPuzzle) Init() tea.Cmd {
	return nil
}

func (EmptyPuzzle) View() string {
	return lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Render("Empty puzzle \n\n press q to quit or ctrl+c to exit")
}

func (p EmptyPuzzle) Update(msg tea.Msg) (Puzzle, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q", "ctrl+c":
			return p, tea.Quit
		}
	}
	return p, nil
}

func (EmptyPuzzle) Preview() string {
	return components.PlainPuzzlePreview("Empty puzzle \n\n press enter or some shit\n yada yada ya")
}
