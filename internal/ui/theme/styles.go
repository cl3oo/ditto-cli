package theme

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/config"
)

type Theme struct {
	Colors ThemeColors

	// Semantic Styles
	Header                 lipgloss.Style
	Title                  lipgloss.Style
	Post                   lipgloss.Style
	Selected               lipgloss.Style
	SelectedCard           lipgloss.Style
	SelectedRow            lipgloss.Style
	Error                  lipgloss.Style
	ErrorBanner            lipgloss.Style
	Text                   lipgloss.Style
	TextSubtle             lipgloss.Style
	AccentText             lipgloss.Style
	Success                lipgloss.Style
	SuccessBanner          lipgloss.Style
	FooterSurface          lipgloss.Style
	FooterSectionTitle     lipgloss.Style
	FooterSectionTitleAlt  lipgloss.Style
	FooterItem             lipgloss.Style
	FooterItemStrong       lipgloss.Style
	CommandSurface         lipgloss.Style
	Surface                lipgloss.Style
	Modal                  lipgloss.Style
	ModalCritical          lipgloss.Style
	Button                 lipgloss.Style
	ButtonFocused          lipgloss.Style
	ButtonSecondary        lipgloss.Style
	ButtonSecondaryFocused lipgloss.Style

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
	background := lipgloss.Color(colors.Background)
	accentContrast := lipgloss.Color(contrastColor(colors.Accent))
	backgroundContrast := lipgloss.Color(contrastColor(colors.Background))

	return Theme{
		Colors: colors,
		Accent: accent,

		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			MarginBottom(1),

		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(accentContrast).
			Background(accent).
			Padding(0, 1),

		Post: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border).
			Padding(0, 1).
			MarginBottom(1),

		Selected: lipgloss.NewStyle().
			Foreground(accent).
			Bold(true),

		SelectedCard: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accent).
			Padding(0, 1).
			MarginBottom(1),

		SelectedRow: lipgloss.NewStyle().
			Background(backgroundContrast).
			PaddingLeft(2),

		Error: lipgloss.NewStyle().
			Foreground(errColor).
			Bold(true),

		ErrorBanner: lipgloss.NewStyle().
			Foreground(errColor).
			Bold(true).
			Padding(0, 1),

		Success: lipgloss.NewStyle().
			Foreground(success).
			Bold(true),

		SuccessBanner: lipgloss.NewStyle().
			Foreground(success).
			Bold(true).
			Padding(0, 1),

		Text: lipgloss.NewStyle().
			Foreground(fg),

		TextSubtle: lipgloss.NewStyle().
			Foreground(subtle),

		AccentText: lipgloss.NewStyle().
			Foreground(accent),

		FooterSurface: lipgloss.NewStyle().
			Background(backgroundContrast).
			Width(0),

		FooterSectionTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(accentContrast).
			Background(accent).
			Padding(0, 1),

		FooterSectionTitleAlt: lipgloss.NewStyle().
			Bold(true).
			Foreground(accent).
			Padding(0, 1),

		FooterItem: lipgloss.NewStyle().
			Foreground(fg).
			Padding(0, 1),

		FooterItemStrong: lipgloss.NewStyle().
			Foreground(fg).
			Bold(true).
			Padding(0, 1),

		CommandSurface: lipgloss.NewStyle().
			Foreground(accentContrast).
			Background(accent).
			Padding(0, 1),

		Surface: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border).
			Padding(1, 2),

		Modal: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accent).
			Padding(1, 2),

		ModalCritical: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errColor).
			Padding(1, 2),

		Button: lipgloss.NewStyle().
			Foreground(fg),

		ButtonFocused: lipgloss.NewStyle().
			Foreground(accentContrast).
			Background(accent).
			Bold(true),

		ButtonSecondary: lipgloss.NewStyle().
			Foreground(subtle),

		ButtonSecondaryFocused: lipgloss.NewStyle().
			Foreground(background).
			Background(subtle).
			Bold(true),

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

func (t Theme) RenderPrimaryButton(label string, focused bool) string {
	if focused {
		return t.ButtonFocused.Render(label)
	}
	return t.Button.Render(label)
}

func (t Theme) RenderSecondaryButton(label string, focused bool) string {
	if focused {
		return t.ButtonSecondaryFocused.Render(label)
	}
	return t.ButtonSecondary.Render(label)
}

func contrastColor(hex string) string {
	hex = strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(hex) != 6 {
		return "#FAFAFA"
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return "#FAFAFA"
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return "#FAFAFA"
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return "#FAFAFA"
	}

	brightness := ((299 * r) + (587 * g) + (114 * b)) / 1000
	if brightness >= 140 {
		return "#111111"
	}
	return "#FAFAFA"
}

func (t Theme) DebugString() string {
	return fmt.Sprintf("accent=%s bg=%s fg=%s", t.Colors.Accent, t.Colors.Background, t.Colors.Foreground)
}
