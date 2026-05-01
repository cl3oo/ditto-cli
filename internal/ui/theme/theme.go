package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/config"
)

type ThemeColors struct {
	Accent     string
	Background string
	Foreground string
	Subtle     string
	Border     string
	Success    string
	Error      string
}

var (
	TokyoNight = ThemeColors{
		Accent:     "#7D56F4",
		Background: "#1A1B26",
		Foreground: "#A9B1D6",
		Subtle:     "#565F89",
		Border:     "#874BFD",
		Success:    "#9ECE6A",
		Error:      "#F7768E",
	}

	Dracula = ThemeColors{
		Accent:     "#BD93F9",
		Background: "#282A36",
		Foreground: "#F8F8F2",
		Subtle:     "#6272A4",
		Border:     "#BD93F9",
		Success:    "#50FA7B",
		Error:      "#FF5555",
	}

	Nord = ThemeColors{
		Accent:     "#88C0D0",
		Background: "#2E3440",
		Foreground: "#D8DEE9",
		Subtle:     "#4C566A",
		Border:     "#81A1C1",
		Success:    "#A3BE8C",
		Error:      "#BF616A",
	}

	Catppuccin = ThemeColors{
		Accent:     "#CBA6F7",
		Background: "#1E1E2E",
		Foreground: "#CDD6F4",
		Subtle:     "#585B70",
		Border:     "#B4BEFE",
		Success:    "#A6E3A1",
		Error:      "#F38BA8",
	}

	// Default Dark Fallback
	DarkTheme = TokyoNight

	// Default Light Fallback
	LightTheme = ThemeColors{
		Accent:     "#7D56F4",
		Background: "#FAFAFA",
		Foreground: "#333333",
		Subtle:     "#999999",
		Border:     "#7D56F4",
		Success:    "#01BE85",
		Error:      "#FF0000",
	}
)

func ResolveThemeColors(cfg config.AppearanceConfig) ThemeColors {
	var theme ThemeColors

	switch cfg.Theme {
	case "tokyo-night":
		theme = TokyoNight
	case "dracula":
		theme = Dracula
	case "nord":
		theme = Nord
	case "catppuccin":
		theme = Catppuccin
	case "light":
		theme = LightTheme
	case "dark":
		theme = DarkTheme
	default:
		if lipgloss.HasDarkBackground() {
			theme = DarkTheme
		} else {
			theme = LightTheme
		}
	}

	// Apply overrides
	if cfg.AccentColor != "" {
		theme.Accent = cfg.AccentColor
	}
	if cfg.BackgroundColor != "" {
		theme.Background = cfg.BackgroundColor
	}
	if cfg.BorderColor != "" {
		theme.Border = cfg.BorderColor
	}
	if cfg.SuccessColor != "" {
		theme.Success = cfg.SuccessColor
	}
	if cfg.ErrorColor != "" {
		theme.Error = cfg.ErrorColor
	}

	return theme
}
