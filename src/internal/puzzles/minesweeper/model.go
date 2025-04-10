package minesweeper

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
)

type Minesweeper struct {
	keys         keyMap
	stopwatch    stopwatch.Model
	prefs        preferences
	minefield    [][]cell
	cursorX      int
	cursorY      int
	isGameOver   bool
	isRunning    bool
	screenHeight int
	screenWidth  int
}

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

type preferences struct {
	width         int
	height        int
	numberOfMines int
	isDebug       bool
	showHelp      bool
}

type cell struct {
	isMine     bool
	isFlagged  bool
	isRevealed bool
}

type point struct {
	x int
	y int
}
