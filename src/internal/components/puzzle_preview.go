package components

import (
	"github.com/Liuuner/go-puzzles/src/internal/common"
	"github.com/charmbracelet/lipgloss"
)

func PlainPuzzlePreview(title string) string {
	return lipgloss.NewStyle().
		Width(common.Config.SelectionContainerWidth).
		Height(common.Config.SelectionContainerHeight).
		Align(lipgloss.Center, lipgloss.Center).
		Render(title)
}

func SimplePuzzlePreview(title string, icon string) string {
	iconHeight := lipgloss.Height(icon)
	iconWidth := lipgloss.Width(icon)
	if iconHeight >= common.Config.SelectionContainerHeight || iconWidth > common.Config.SelectionContainerWidth {
		return PlainPuzzlePreview(lipgloss.NewStyle().Foreground(lipgloss.Color("#F00")).Render("Icon for " + title + " is too big"))
	}

	topPart := lipgloss.Place(common.Config.SelectionContainerWidth, iconHeight, lipgloss.Center, lipgloss.Top, icon)
	titlePart := lipgloss.Place(common.Config.SelectionContainerWidth, common.Config.SelectionContainerHeight-iconHeight, lipgloss.Center, lipgloss.Top, lipgloss.NewStyle().Bold(true).Render(title))
	return lipgloss.JoinVertical(lipgloss.Center, topPart, titlePart)
}
