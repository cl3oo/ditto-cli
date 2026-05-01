package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/config"
)

type Theme struct {
	Colors ThemeColors

	// Semantic Styles
	Header     lipgloss.Style
	Title      lipgloss.Style
	Post       lipgloss.Style
	Selected   lipgloss.Style
	Error      lipgloss.Style
	Text       lipgloss.Style
	TextSubtle lipgloss.Style
	AccentText lipgloss.Style
	Success    lipgloss.Style

	Accent   lipgloss.Color
	Markdown string // Markdown style json
}

func NewTheme(cfg config.AppearanceConfig) Theme {
	colors := ResolveThemeColors(cfg)

	accent := lipgloss.Color(colors.Accent)
	border := lipgloss.Color(colors.Border)
	success := lipgloss.Color(colors.Success)
	errColor := lipgloss.Color(colors.Error)
	fg := lipgloss.Color(colors.Foreground)
	subtle := lipgloss.Color(colors.Subtle)

	return Theme{
		Colors: colors,
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

		Success: lipgloss.NewStyle().
			Foreground(success).
			Bold(true),

		Text: lipgloss.NewStyle().
			Foreground(fg),

		TextSubtle: lipgloss.NewStyle().
			Foreground(subtle),

		AccentText: lipgloss.NewStyle().
			Foreground(accent),

		Markdown: `{
			"para": { "margin_left": 2 },
			"heading": { "color": "` + colors.Accent + `", "bold": true },
			"link": { "color": "` + colors.Accent + `", "underline": true },
			"code": { "background_color": "` + colors.Background + `" }
		}`,
	}
}

func DefaultTheme() Theme {
	return NewTheme(config.DefaultConfig().Appearance)
}
