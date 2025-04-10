package cmd

import (
	"github.com/Liuuner/go-puzzles/src/internal/puzzles"
)

type model struct {
	terminalInfo   TerminalInfo
	layout         layout
	puzzleOpened   bool
	selectedPuzzle int
	puzzle         puzzles.Puzzle
	puzzles        []puzzles.Puzzle
	standaloneMode bool
}

type layout struct {
	headerHeight        int
	selectionHeight     int
	selectionItemsInRow int
}

type TerminalInfo struct {
	fullWidth  int
	fullHeight int
}
