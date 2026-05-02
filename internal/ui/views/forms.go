package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type formSection struct {
	Label   string
	Content string
}

func renderFormShell(t theme.Theme, viewportWidth int, title, subtitle string, sections []formSection, actions []string, feedback string) string {
	contentWidth := formContentWidth(viewportWidth)
	lines := []string{t.Title.Render(" " + title + " ")}
	if subtitle != "" {
		lines = append(lines, "", t.TextSubtle.Render(subtitle))
	}

	for _, section := range sections {
		if section.Content == "" {
			continue
		}
		lines = append(lines, "")
		if section.Label != "" {
			lines = append(lines, t.Text.Bold(true).Render(section.Label))
		}
		lines = append(lines, section.Content)
	}

	if len(actions) > 0 {
		lines = append(lines, "", renderFormActions(contentWidth, actions))
	}

	if feedback != "" {
		lines = append(lines, "", feedback)
	}

	boxWidth := formBoxWidth(viewportWidth)
	box := t.Surface.Copy().Width(boxWidth)
	rendered := box.Render(strings.Join(lines, "\n"))
	if viewportWidth > 0 {
		return lipgloss.Place(viewportWidth, lipgloss.Height(rendered), lipgloss.Center, lipgloss.Top, rendered)
	}
	return rendered
}

func renderFormActions(contentWidth int, actions []string) string {
	joined := strings.Join(actions, "  ")
	if contentWidth <= 0 || lipgloss.Width(joined) <= contentWidth {
		return joined
	}
	return strings.Join(actions, "\n")
}

func renderFormFeedback(t theme.Theme, tone, title, body string) string {
	return RenderStateSurface(t, 0, tone, title, body, nil)
}

func formBoxWidth(viewportWidth int) int {
	if viewportWidth <= 0 {
		return 72
	}
	return min(max(58, viewportWidth-12), 88)
}

func formContentWidth(viewportWidth int) int {
	return formBoxWidth(viewportWidth) - 6
}

func fitInputWidth(viewportWidth int, fallback int) int {
	width := formContentWidth(viewportWidth)
	if width <= 0 {
		return fallback
	}
	return max(20, width)
}
