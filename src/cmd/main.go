package cmd

import (
	"flag"
	"fmt"
	"github.com/Liuuner/go-puzzles/src/internal/common"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	DEFAULT_WIDTH  = 9
	DEFAULT_HEIGHT = 9
	DEFAULT_MINES  = 10
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Sweep  key.Binding
	Flag   key.Binding
	New    key.Binding
	Redraw key.Binding
	Help   key.Binding
	Quit   key.Binding
}

var DefaultKeyMap = keyMap{
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
	Sweep: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "sweep"),
	),
	Flag: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "flag"),
	),
	New: key.NewBinding(
		key.WithKeys("N"),
		key.WithHelp("N", "new game"),
	),
	Redraw: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "redraw"),
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

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},    // first column
		{k.Sweep, k.Flag, k.New, k.Redraw}, // first column
		{k.Help, k.Quit},                   // second column
	}
}

func Run() {
	var wFlag, hFlag, numMinesFlag int
	var previewTheme bool

	flag.IntVar(&wFlag, "w", DEFAULT_WIDTH, "minefield width")
	flag.IntVar(&hFlag, "h", DEFAULT_HEIGHT, "minefield height")
	flag.IntVar(&numMinesFlag, "n", DEFAULT_MINES, "number of mines")
	flag.BoolVar(&previewTheme, "preview", false, "preview theme")

	flag.Parse()

	prefs := preferences{
		width:         wFlag,
		height:        hFlag,
		numberOfMines: numMinesFlag,
		showHelp:      true,
		isDebug:       false,
	}

	m := initialModel(prefs)

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}

func testSolver() {
	prefs := preferences{
		width:         DEFAULT_WIDTH,
		height:        DEFAULT_HEIGHT,
		numberOfMines: DEFAULT_MINES,
		showHelp:      true,
		isDebug:       false,
	}
	tries := 0

	for {
		tries++
		m := initialModel(prefs)
		placeMinesSkipCursor(m.minefield, prefs, [2]int{m.cursorX, m.cursorY})
		sweep(m.cursorX, m.cursorY, &m, true, make(common.Set[point]))
		sb1 := strings.Builder{}
		sb1.WriteString("Trying to solve:\n")
		writeAsciiMinefield(&sb1, m)

		sb2 := strings.Builder{}
		ok := solveMinesweeper(&m)

		if ok {
			sb2.WriteString("Solved:\n")
		} else {
			sb2.WriteString(fmt.Sprintf("Failed to solve: Remaining: %d\n", minesLeft(m)))
		}
		writeAsciiMinefield(&sb2, m)

		sb3 := strings.Builder{}
		if !ok {
			// reveal cells
			for y := range m.minefield {
				for x := range m.minefield[y] {
					cell := &m.minefield[y][x]
					if !cell.isRevealed {
						cell.isRevealed = true
					}
				}
			}
			sb3.WriteString("Revealed:\n")
			writeAsciiMinefield(&sb3, m)
			sb3.WriteString("\n")
		}

		println(lipgloss.JoinHorizontal(0, sb1.String(), sb2.String(), sb3.String()))

		if ok {
			break
		}
	}

	fmt.Printf("Solved with %d tries\n", tries)
}

func placeMinesSkipCursor(minefield [][]cell, prefs preferences, cursorPos [2]int) {
	if prefs.numberOfMines >= prefs.width*prefs.height {
		panic("Too many mines for the given field size")
	}

	positions := make([][2]int, 0, prefs.width*prefs.height)
	for y := 0; y < prefs.height; y++ {
		for x := 0; x < prefs.width; x++ {
			// don't place mines on the cursor
			if x == cursorPos[0] && y == cursorPos[1] {
				continue
			}
			positions = append(positions, [2]int{x, y})
		}
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(positions), func(i, j int) { positions[i], positions[j] = positions[j], positions[i] })

	for i := 0; i < prefs.numberOfMines; i++ {
		x, y := positions[i][0], positions[i][1]
		minefield[y][x].isMine = true
	}
}

func debugPrintMinefield(minefield [][]cell) {
	s := strings.Builder{}
	for y := range minefield {
		for _, c := range minefield[y] {
			if c.isMine {
				s.WriteString("M")
			} else {
				s.WriteString("-")
			}
		}
		s.WriteString("\n")
	}
	fmt.Println(s.String())
}

func createEmptyMinefield(prefs preferences) [][]cell {
	minefield := make([][]cell, prefs.height)

	for y := range minefield {
		minefield[y] = make([]cell, prefs.width)
		for x := range minefield[y] {
			minefield[y][x] = cell{}
		}
	}
	return minefield
}

func initialModel(prefs preferences) model {
	minefield := createEmptyMinefield(prefs)

	//debugPrintMinefield(minefield)

	// TODO instantiate the mines after the first sweep to make sure first click never hits a mine
	//placeMines(minefield, prefs)

	//debugPrintMinefield(minefield)

	return model{
		stopwatch: stopwatch.NewWithInterval(time.Second),
		keys:      DefaultKeyMap,
		prefs:     prefs,
		minefield: minefield,
		cursorX:   prefs.width / 2,
		cursorY:   prefs.height / 2,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.stopwatch.Init(),
		m.stopwatch.Stop(),
		m.stopwatch.Reset(),
		tea.SetWindowTitle("Go-Sweep"),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cursorCell := &m.minefield[m.cursorY][m.cursorX]

	var cmd tea.Cmd
	if m.isRunning && !m.isGameOver {
		m.stopwatch, cmd = m.stopwatch.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width
		m.screenHeight = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Up):
			m.cursorY--
			if m.cursorY < 0 {
				m.cursorY = m.prefs.height - 1
			}
		case key.Matches(msg, DefaultKeyMap.Down):
			m.cursorY++
			if m.cursorY > m.prefs.height-1 {
				m.cursorY = 0
			}
		case key.Matches(msg, DefaultKeyMap.Left):
			m.cursorX--
			if m.cursorX < 0 {
				m.cursorX = m.prefs.width - 1
			}
		case key.Matches(msg, DefaultKeyMap.Right):
			m.cursorX++
			if m.cursorX > m.prefs.width-1 {
				m.cursorX = 0
			}
		case key.Matches(msg, DefaultKeyMap.New):
			isDebug := m.prefs.isDebug
			//y, x := m.cursorY, m.cursorX
			showHelp := m.prefs.showHelp
			m = initialModel(m.prefs)
			m.isRunning = false
			m.prefs.isDebug = isDebug
			//m.cursorY, m.cursorX = y, x
			m.prefs.showHelp = showHelp
			break
		case key.Matches(msg, DefaultKeyMap.Sweep):
			if m.isGameOver {
				break
			}

			if !m.isRunning {
				m.isRunning = true
				placeMinesSkipCursor(m.minefield, m.prefs, [2]int{m.cursorX, m.cursorY})
				cmd = m.stopwatch.Start()
			}
			sweep(m.cursorX, m.cursorY, &m, true, make(common.Set[point]))

			if checkDidWin(m) {
				m.isGameOver = true
			}
		case key.Matches(msg, DefaultKeyMap.Flag):
			if !m.isRunning && !m.isGameOver {
				break
			}
			if m.isGameOver {
				break
			}
			if cursorCell.isRevealed {
				sweep(m.cursorX, m.cursorY, &m, true, make(common.Set[point]))
			} else {
				cursorCell.isFlagged = !cursorCell.isFlagged
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case msg.String() == "S":
			solveMinesweeper(&m)
		}
	}

	return m, cmd
}

func (m model) View() string {
	var sb strings.Builder
	writeHeader(&sb, m)
	sb.WriteString("\n\n")

	writeAsciiMinefield(&sb, m)

	sb.WriteString("\n\n")
	writeHelp(&sb, m)
	return sb.String()
}

func writeHeader(sb *strings.Builder, m model) {
	if m.isGameOver {
		sb.WriteString("Game Over! \n")
		if checkDidWin(m) {
			sb.WriteString("You WON!!\n")
		} else {
			sb.WriteString("You lost...\n")
		}
	} else {
		sb.WriteString("...go sweep...\n")
		sb.WriteString(fmt.Sprintf("%v mines left\n", minesLeft(m)))
	}
	sb.WriteString(m.stopwatch.View())
	sb.WriteString(" elapsed\n")
}

func writeAsciiMinefield(sb *strings.Builder, m model) {
	cursorChar := "⯀"
	mineCharacter := "B"
	flagCharacter := "▶"
	//mineCharacter := "⬤"

	strs := make([][]string, m.prefs.height)
	for y, row := range m.minefield {
		strs[y] = make([]string, m.prefs.width)
		for x, c := range row {
			switch {
			case (m.isGameOver || m.prefs.isDebug) && c.isMine:
				strs[y][x] = mineCharacter
			case c.isRevealed:
				strs[y][x] = asciiViewForMineAtPosition(x, y, m)
			case c.isFlagged:
				strs[y][x] = flagCharacter
			default:
				strs[y][x] = " "
			}
		}
	}

	unrevieldColor := "#313244"

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(strs...).
		StyleFunc(func(row, col int) lipgloss.Style {
			var fg lipgloss.TerminalColor = lipgloss.NoColor{}
			var bg lipgloss.TerminalColor = lipgloss.NoColor{}

			s := lipgloss.NewStyle().
				Padding(0, 1)

			switch strs[row][col] {
			case "0":
				strs[row][col] = ""
			case "1":
				fg = lipgloss.Color("#74adf2")
			case "2":
				fg = lipgloss.Color("#00FF00")
			case "3":
				fg = lipgloss.Color("#FF0000")
			case "4":
				fg = lipgloss.Color("#28706d")
			case "5":
				fg = lipgloss.Color("#b06446")
			case "6":
				fg = lipgloss.Color("#FF0000")
			case "7":
				fg = lipgloss.Color("#8a7101")
			case "8":
				fg = lipgloss.Color("#111")
				bg = lipgloss.Color("#bfbfbf")
			/*case "*":
			fg = lipgloss.Color("#FF33FF")*/
			case flagCharacter:
				fg = lipgloss.Color("#ffee00")
				//bg = lipgloss.Color(unrevieldColor)

				//fg = lipgloss.Color("#111")
			case mineCharacter:
				fg = lipgloss.Color("#FF0000")
				s = s.Bold(true)
				//fg = lipgloss.Color("#111")
			case " ", cursorChar:
				bg = lipgloss.Color(unrevieldColor)
			}

			// if the cursor is on this cell, highlight it
			if m.cursorX == col && m.cursorY == row {
				if strs[row][col] == " " || strs[row][col] == cursorChar {
					strs[row][col] = cursorChar
					//bg = lipgloss.Color("#313244")
					fg = lipgloss.Color("#7f849c")

					//strs[row][col] = "█"
					//fg = lipgloss.Color("#bfbfbf")
					//bg = lipgloss.Color("#bfbfbf")
				} else {
					noCol := lipgloss.NoColor{}
					if bg != noCol {
						fg = bg
					}
					bg = lipgloss.Color("#585b70")
					//bg = lipgloss.Color("#313244")
				}
			}

			return s.
				Foreground(fg).
				Background(bg)
		})

	sb.WriteString(t.Render())
}

func writeHelp(sb *strings.Builder, m model) {
	helpView := m.help.View(m.keys)
	// prepare the string builder so that help is at bottom of papge
	calculatedHeight := m.screenHeight - strings.Count(sb.String(), "\n") - strings.Count(helpView, "\n")
	paddingHeight := math.Max(float64(calculatedHeight), 0)
	sb.WriteString(strings.Repeat("\n", int(paddingHeight)))
	sb.WriteString(helpView)
	//sb.WriteString("\n")
}

func sweep(x, y int, m *model, userInitiatedSweep bool, swept common.Set[point]) {
	cell := &m.minefield[y][x]

	if cell.isRevealed && userInitiatedSweep {
		adjMines := countAdjacentMines(x, y, *m)
		adjFlags := countAdjacentFlags(x, y, *m)
		if adjFlags >= adjMines {
			autoSweep(x, y, m)
		}
		return
	}

	if cell.isMine {
		if userInitiatedSweep {
			m.isGameOver = true
		}
		return
	}

	touching := countAdjacentMines(x, y, *m)

	p := point{x: x, y: y}
	if touching == 0 && !swept.Has(p) {
		swept.Add(p)
		forEachSurroundingCellDo(x, y, m, func(x, y int, m *model) {
			sweep(x, y, m, false, swept)
		})
	}

	cell.isRevealed = true
}

func minesLeft(m model) int {
	flags := 0
	for y := range m.minefield {
		for _, mine := range m.minefield[y] {
			if mine.isFlagged && !mine.isRevealed {
				flags++
			}
		}
	}
	return m.prefs.numberOfMines - flags
}

func checkDidWin(m model) bool {
	for y := range m.minefield {
		for _, mine := range m.minefield[y] {
			if !mine.isMine && !mine.isRevealed {
				return false
			}
		}
	}
	return true
}

func asciiViewForMineAtPosition(x, y int, m model) string {
	if m.minefield[y][x].isMine {
		return "B"
	}
	return fmt.Sprint(countAdjacentMines(x, y, m))
}

func forEachSurroundingCellDo(x, y int, m *model, do func(x, y int, m *model)) {
	w := m.prefs.width
	h := m.prefs.height
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if (dx == 0 && dy == 0) || x+dx < 0 || x+dx > w-1 || y+dy < 0 || y+dy > h-1 {
				continue
			}
			do(x+dx, y+dy, m)
		}
	}
}

func autoSweep(x, y int, m *model) {
	forEachSurroundingCellDo(x, y, m, func(x, y int, m *model) {
		cell := m.minefield[y][x]
		if !cell.isRevealed && !cell.isFlagged {
			sweep(x, y, m, true, make(common.Set[point]))
		}
	})
}

func countAdjacentFlags(x, y int, m model) int {
	adj := 0
	forEachSurroundingCellDo(x, y, &m, func(x, y int, m *model) {
		if m.minefield[y][x].isFlagged {
			adj++
		}
	})
	return adj
}

func countAdjacentMines(x, y int, m model) int {
	adj := 0
	forEachSurroundingCellDo(x, y, &m, func(x, y int, m *model) {
		if m.minefield[y][x].isMine {
			adj++
		}
	})
	return adj
}

func isSolvable(m *model) bool {
	mCopy := initialModel(m.prefs)
	duplicate := make([][]cell, len(m.minefield))
	for i := range m.minefield {
		duplicate[i] = make([]cell, len(m.minefield[i]))
		copy(duplicate[i], m.minefield[i])
	}
	mCopy.minefield = duplicate

	return solveMinesweeper(&mCopy)
}

/* ----------------------------------------------------------------------
* Minesweeper solver, used to ensure the generated grids are
* solvable without having to take risks.
 */
func solveMinesweeper(m *model) bool {
	for {
		progress := false

		for y := range m.minefield {
			for x := range m.minefield[y] {
				cell := &m.minefield[y][x]
				if cell.isRevealed {
					adjMines := countAdjacentMines(x, y, *m)
					adjFlags := countAdjacentFlags(x, y, *m)
					adjHidden := countAdjacentHidden(x, y, *m)

					// Flag all adjacent hidden cells if the number of adjacent mines equals the number of adjacent flags plus hidden cells
					if adjMines == adjFlags+adjHidden {
						forEachSurroundingCellDo(x, y, m, func(x, y int, m *model) {
							if !m.minefield[y][x].isRevealed && !m.minefield[y][x].isFlagged {
								m.minefield[y][x].isFlagged = true
								progress = true
							}
						})
					}

					// Reveal all adjacent hidden cells if the number of adjacent flags equals the number of adjacent mines
					if adjMines == adjFlags {
						forEachSurroundingCellDo(x, y, m, func(x, y int, m *model) {
							if !m.minefield[y][x].isRevealed && !m.minefield[y][x].isFlagged {
								sweep(x, y, m, false, make(common.Set[point]))
								progress = true
							}
						})
					}
				}
			}
		}

		// If no progress was made, break the loop
		if !progress {
			break
		}
	}
	// check if there are any unrevealed cells left
	for y := range m.minefield {
		for x := range m.minefield[y] {
			if !m.minefield[y][x].isRevealed && !m.minefield[y][x].isFlagged {
				return false
			}
		}
	}
	return true
}

func countAdjacentHidden(x, y int, m model) int {
	adj := 0
	forEachSurroundingCellDo(x, y, &m, func(x, y int, m *model) {
		if !m.minefield[y][x].isRevealed && !m.minefield[y][x].isFlagged {
			adj++
		}
	})
	return adj
}
