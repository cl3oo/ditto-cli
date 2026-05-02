package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

func renderVerticalCard(cardStyle lipgloss.Style, isSelected bool, t theme.Theme, meta, title, body, stats string) string {
	titleStyle := t.Text.Bold(true)
	if isSelected {
		titleStyle = t.AccentText.Bold(true)
	}

	lines := []string{t.TextSubtle.Render(meta), titleStyle.Render(title)}
	if body != "" {
		lines = append(lines, t.Text.Render(body))
	}
	if stats != "" {
		lines = append(lines, t.TextSubtle.Render(stats))
	}

	return withSelectionIndicator(cardStyle.Render(strings.Join(lines, "\n\n")), isSelected, t)
}

func clampLine(text string, limit int) string {
	text = strings.TrimSpace(strings.Join(strings.Fields(text), " "))
	if text == "" || limit <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	if limit == 1 {
		return "…"
	}
	return strings.TrimSpace(string(runes[:limit-1])) + "…"
}
