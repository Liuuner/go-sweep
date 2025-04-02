package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/stopwatch"
)

type model struct {
	keys         keyMap
	help         help.Model
	screenHeight int
	stopwatch    stopwatch.Model
	prefs        preferences
	minefield    [][]cell
	cursorX      int
	cursorY      int
	isGameOver   bool
	isRunning    bool
}

type preferences struct {
	width          int
	height         int
	numberOfMines  int
	isDebug        bool
	showHelp       bool
	shouldUseEmoji bool
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
