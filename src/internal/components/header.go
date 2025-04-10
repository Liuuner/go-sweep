package components

import "github.com/charmbracelet/lipgloss"

func Header(width int, start, title, end string) string {
	// Calculate the available space for the header
	if width < len(start+title+end) {
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			Foreground(lipgloss.Color("#F00")).
			Render(lipgloss.PlaceHorizontal(width, lipgloss.Center, "Header too long"))
	}

	header := lipgloss.PlaceHorizontal(width, lipgloss.Center, title)
	header = overwriteStart(header, start)
	header = overwriteEnd(header, end)

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

// overwriteStart overwrites the start of the base string with the overwrite string
func overwriteStart(base, overwrite string) string {
	baseRunes := []rune(base)
	overwriteRunes := []rune(overwrite)

	if len(overwriteRunes) > len(baseRunes) {
		// If overwrite string is longer than base, return overwrite
		return overwrite
	}

	// Overwrite the first characters
	copy(baseRunes[:len(overwriteRunes)], overwriteRunes)
	return string(baseRunes)
}
