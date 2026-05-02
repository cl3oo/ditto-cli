package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

func RenderCenteredStateSurface(t theme.Theme, viewportWidth, viewportHeight int, tone, title, body string, actions []string) string {
	panel := renderStateSurfacePanel(t, viewportWidth, tone, title, body, actions)
	if viewportWidth <= 0 && viewportHeight <= 0 {
		return panel
	}

	return lipgloss.Place(
		max(0, viewportWidth),
		max(0, viewportHeight),
		lipgloss.Center,
		lipgloss.Center,
		panel,
	)
}

func RenderStateSurface(t theme.Theme, viewportWidth int, tone, title, body string, actions []string) string {
	return renderStateSurfacePanel(t, viewportWidth, tone, title, body, actions)
}

func renderStateSurfacePanel(t theme.Theme, viewportWidth int, tone, title, body string, actions []string) string {
	box, heading, bodyStyle := stateSurfaceStyles(t, tone)
	lines := make([]string, 0, 3+len(actions))

	if title != "" {
		lines = append(lines, heading.Render(title))
	}
	if body != "" {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, bodyStyle.Render(body))
	}
	if len(actions) > 0 {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, t.Text.Render("Try:"))
		for _, action := range actions {
			lines = append(lines, "  • "+t.TextSubtle.Render(action))
		}
	}

	if viewportWidth > 0 {
		width := min(max(52, viewportWidth-12), 88)
		return box.Width(width).Render(strings.Join(lines, "\n"))
	}

	return box.Render(strings.Join(lines, "\n"))
}

func stateSurfaceStyles(t theme.Theme, tone string) (lipgloss.Style, lipgloss.Style, lipgloss.Style) {
	box := t.Surface.Copy().BorderForeground(t.TextSubtle.GetForeground())
	heading := t.Title.Copy().Background(lipgloss.Color(t.Colors.Background)).Foreground(t.Accent).Padding(0, 0)
	body := t.Text

	switch tone {
	case "loading":
		box = t.Surface.Copy().BorderForeground(t.Accent)
		heading = t.AccentText.Bold(true)
	case "success":
		box = t.Surface.Copy().BorderForeground(t.Success.GetForeground())
		heading = t.Success
	case "error":
		box = t.Surface.Copy().BorderForeground(t.Error.GetForeground())
		heading = t.Error
	case "empty":
		fallthrough
	default:
	}

	return box, heading, body
}
