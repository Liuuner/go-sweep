package cmd

import (
	"fmt"
	"github.com/Liuuner/go-puzzles/src/internal/common"
	"github.com/Liuuner/go-puzzles/src/internal/puzzles"
	"github.com/Liuuner/go-puzzles/src/internal/puzzles/minesweeper"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

func Run() {
	//puzzle := minesweeper.Minesweeper{}.New()
	puzzle := puzzles.EmptyPuzzle{}

	m := model{
		puzzle:         puzzle,
		puzzleOpened:   false,
		selectedPuzzle: 0,
		puzzles:        []puzzles.Puzzle{minesweeper.Minesweeper{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}},
		layout: layout{
			selectionItemsInRow: 1,
		},
	}

	p := tea.NewProgram(m, tea.WithAltScreen() /*, tea.WithMouseCellMotion()*/)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Aye! There's been an error: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.handleWindowResize(msg)
	}

	if _, ok := msg.(common.QuitGameMsg); ok {
		m.puzzleOpened = false
		m.puzzle = puzzles.EmptyPuzzle{}
	}

	var batch tea.Cmd

	if m.puzzleOpened {
		puzzle, cmd := m.puzzle.Update(msg)
		m.puzzle = puzzle
		batch = tea.Batch(batch, cmd)
	} else {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			m, batch = m.handleKeyMsg(keyMsg)
			if m.puzzleOpened {
				batch = tea.Batch(batch, m.puzzle.Init())
			}
		}
	}

	return m, batch
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (model, tea.Cmd) {
	switch {
	case key.Matches(msg, common.Hotkeys.Quit):
		return m, tea.Quit
	case key.Matches(msg, common.Hotkeys.Up):
		newSelected := m.selectedPuzzle - m.layout.selectionItemsInRow
		if newSelected >= 0 {
			m.selectedPuzzle = newSelected
		}
	case key.Matches(msg, common.Hotkeys.Down):
		newSelected := m.selectedPuzzle + m.layout.selectionItemsInRow
		if newSelected < len(m.puzzles) {
			m.selectedPuzzle = newSelected
		}
	case key.Matches(msg, common.Hotkeys.Left):
		if m.selectedPuzzle%3 != 0 {
			m.selectedPuzzle = max(0, m.selectedPuzzle-1)
		}
	case key.Matches(msg, common.Hotkeys.Right):
		if (m.selectedPuzzle+1)%m.layout.selectionItemsInRow != 0 {
			m.selectedPuzzle = min(len(m.puzzles)-1, m.selectedPuzzle+1)
		}
	case key.Matches(msg, common.Hotkeys.Select):
		if !m.puzzleOpened {
			m.puzzle = m.puzzles[m.selectedPuzzle].New()
			m.puzzleOpened = true
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.puzzleOpened {
		return m.puzzle.View()
	}

	header := drawHeader(m)

	puzzleSelection := drawPuzzleSelection(m)

	return lipgloss.JoinVertical(lipgloss.Center, header, puzzleSelection)
}

func (m *model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.terminalInfo.fullWidth = msg.Width
	m.terminalInfo.fullHeight = msg.Height

	m.layout.selectionHeight = m.terminalInfo.fullHeight - common.Config.HeaderHeight - 1

	m.layout.selectionItemsInRow = max(m.terminalInfo.fullWidth/(common.Config.SelectionContainerWidth+3), 1)

	if m.puzzleOpened {
		m.puzzle.Update(msg)
	}
}

func drawHeader(m model) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Render("⚄ Terminal Puzzles ⚄")

	help := "Help: ?           "

	header := lipgloss.PlaceHorizontal(m.terminalInfo.fullWidth, lipgloss.Center, title)
	header = overwriteEnd(header, help)

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Render(header)
}

// overwriteEnd overwrites the end of the base string with the overwrite string
func overwriteEnd(base, overwrite string) string {
	baseRunes := []rune(base)
	overwriteRunes := []rune(overwrite)

	if len(overwriteRunes) > len(baseRunes) {
		// If overwrite string is longer than base, return overwrite
		return overwrite
	}

	// Overwrite the last characters
	copy(baseRunes[len(baseRunes)-len(overwriteRunes):], overwriteRunes)
	return string(baseRunes)
}

func drawPuzzleSelection(m model) string {
	sb := strings.Builder{}
	width := (common.Config.SelectionContainerWidth+2)*m.layout.selectionItemsInRow + m.layout.selectionItemsInRow - 1

	lineAmount := len(m.puzzles) / m.layout.selectionItemsInRow
	if len(m.puzzles)%m.layout.selectionItemsInRow != 0 {
		lineAmount++
	}

	grid := make([][]string, lineAmount)
	for i := range grid {
		grid[i] = make([]string, m.layout.selectionItemsInRow)
	}

	for i, p := range m.puzzles {
		selected := i == m.selectedPuzzle

		lineNum := i / m.layout.selectionItemsInRow
		itemNum := i % m.layout.selectionItemsInRow
		if itemNum != 0 {
			grid[lineNum][itemNum] = lipgloss.NewStyle().MarginLeft(1).Render(buildPuzzleContainer(p, selected))
		} else {
			grid[lineNum][itemNum] = buildPuzzleContainer(p, selected)
		}
	}

	for _, line := range grid {
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, line...))
		sb.WriteString("\n")
	}
	return lipgloss.NewStyle().Width(width).Render(sb.String())
}

func buildPuzzleContainer(puzzle puzzles.Puzzle, selected bool) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Height(common.Config.SelectionContainerHeight).
		Width(common.Config.SelectionContainerWidth)

	if selected {
		style = style.BorderForeground(lipgloss.Color("#89b4fa"))
	}

	return style.Render(puzzle.Preview())
}
