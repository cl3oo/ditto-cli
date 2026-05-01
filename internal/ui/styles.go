package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/config"
)

type Theme struct {
	Header   lipgloss.Style
	Title    lipgloss.Style
	Post     lipgloss.Style
	Selected lipgloss.Style
	Error    lipgloss.Style
	Accent   lipgloss.Color
	Markdown string // Markdown style json
}

func NewTheme(cfg config.AppearanceConfig) Theme {
	accent := lipgloss.Color(cfg.AccentColor)
	border := lipgloss.Color(cfg.BorderColor)
	success := lipgloss.Color(cfg.SuccessColor)
	errColor := lipgloss.Color(cfg.ErrorColor)

	return Theme{
		Accent: accent,
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			MarginBottom(1),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(accent).
			Padding(0, 1),

		Post: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border).
			Padding(0, 1).
			MarginBottom(1),

		Selected: lipgloss.NewStyle().
			Foreground(success).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(errColor).
			Bold(true),

		Markdown: `{
			"para": { "margin_left": 2 },
			"heading": { "color": "` + cfg.AccentColor + `", "bold": true },
			"link": { "color": "` + cfg.AccentColor + `", "underline": true },
			"code": { "background_color": "#282a36" }
		}`,
	}
}

// DefaultTheme provides a fallback theme
func DefaultTheme() Theme {
	return NewTheme(config.DefaultConfig().Appearance)
}
