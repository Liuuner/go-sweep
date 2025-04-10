package cmd

import (
	"context"
	"fmt"
	"github.com/Liuuner/go-puzzles/src/internal/common"
	"github.com/Liuuner/go-puzzles/src/internal/components"
	"github.com/Liuuner/go-puzzles/src/internal/puzzles"
	"github.com/Liuuner/go-puzzles/src/internal/puzzles/minesweeper"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"slices"
	"strings"
)

func Run() {

	cmd := &cli.Command{
		EnableShellCompletion: true,
		Name:                  "go-puzzles",
		Version:               "v0.0.1-dev",
		HideVersion:           true,
		HideHelpCommand:       true,
		Description:           "A collection of terminal puzzles",
		Commands: []*cli.Command{
			{
				Name:     "minesweeper",
				Usage:    "Play minesweeper",
				Aliases:  []string{"ms"},
				Category: "Games",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return runStandalone(minesweeper.NewWithCmd(cmd))
				},
				Flags: []cli.Flag{
					cli.HelpFlag,
					&cli.IntFlag{
						Name:        "width",
						Aliases:     []string{"w"},
						Usage:       "Width of the minesweeper field",
						Required:    false,
						HideDefault: true,
					},
					&cli.IntFlag{
						Name:        "height",
						Aliases:     []string{"h"},
						Usage:       "Height of the minesweeper field",
						Required:    false,
						HideDefault: true,
					},
					&cli.IntFlag{
						Name:        "mines",
						Aliases:     []string{"m"},
						Usage:       "Number of mines in the field",
						Required:    false,
						HideDefault: true,
					},
					&cli.StringFlag{
						Name:        "difficulty",
						Usage:       "Difficulty level of the minesweeper game (beginner, intermediate, expert)",
						Aliases:     []string{"d"},
						Required:    false,
						HideDefault: true,
						Validator: func(d string) error {
							difficulties := []string{"beginner", "intermediate", "expert"}
							if !slices.Contains(difficulties, d) {
								return fmt.Errorf("invalid difficulty level: %s, valid levels are: %v", d, difficulties)
							}
							return nil
						},
					},
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runDefault()
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runDefault() error {
	m := model{
		puzzle:         puzzles.EmptyPuzzle{},
		puzzleOpened:   false,
		selectedPuzzle: 0,
		puzzles:        []puzzles.Puzzle{minesweeper.Minesweeper{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}, puzzles.EmptyPuzzle{}},
		layout: layout{
			selectionItemsInRow: 1,
		},
	}

	p := tea.NewProgram(m, tea.WithAltScreen() /*, tea.WithMouseCellMotion()*/)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func runStandalone(puzzle puzzles.Puzzle) error {
	m := model{
		puzzle:         puzzle,
		standaloneMode: true,
	}

	p := tea.NewProgram(m, tea.WithAltScreen() /*, tea.WithMouseCellMotion()*/)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.handleWindowResize(msg)
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	if _, ok := msg.(common.QuitGameMsg); ok {
		if m.standaloneMode {
			return m, tea.Quit
		}
		m.puzzleOpened = false
		m.puzzle = puzzles.EmptyPuzzle{}
	}

	var batch tea.Cmd

	if m.puzzleOpened || m.standaloneMode {
		puzzle, cmd := m.puzzle.Update(msg)
		m.puzzle = puzzle
		batch = tea.Batch(batch, cmd)
	} else {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			// todo handleKeyMsg will tell when a puzzle is selected with a cmd like PuzzleSelected
			m, batch = m.handleKeyMsg(keyMsg)
			if m.puzzleOpened {
				// todo when all m.puzzle.Init() commands are finished, we return the PuzzleReady Cmd or something in this fashion
				batch = tea.Batch(batch, m.puzzle.Init(), tea.WindowSize())
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
	if m.puzzleOpened || m.standaloneMode {
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
	title := lipgloss.NewStyle().Bold(true).Render("Go Puzzles")

	return components.Header(m.terminalInfo.fullWidth, "v0.0.1-dev", title, "Help: ?")
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
