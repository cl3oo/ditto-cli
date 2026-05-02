package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

func selectionPrefix(selected bool, t theme.Theme) string {
	if selected {
		return t.AccentText.Bold(true).Render(">")
	}
	return t.TextSubtle.Render(" ")
}

func withSelectionIndicator(content string, selected bool, t theme.Theme) string {
	prefix := selectionPrefix(selected, t)
	lines := strings.Split(content, "\n")
	for i := range lines {
		linePrefix := " "
		if i == 0 {
			linePrefix = prefix
		}
		lines[i] = lipgloss.JoinHorizontal(lipgloss.Top, linePrefix, " ", lines[i])
	}
	return strings.Join(lines, "\n")
}
