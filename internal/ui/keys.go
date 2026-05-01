package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/rfcku/ditto-cli/internal/config"
)

type KeyMap struct {
	Quit        key.Binding
	Back        key.Binding
	New         key.Binding
	Select      key.Binding
	Upvote      key.Binding
	Downvote    key.Binding
	Refresh     key.Binding
	Login       key.Binding
	Feed        key.Binding
	Communities key.Binding
	Palette     key.Binding
}

func NewKeyMap(cfg config.KeyConfig) KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys(cfg.Quit, "ctrl+c"),
			key.WithHelp(cfg.Quit+"/ctrl+c", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys(cfg.Back, "esc"),
			key.WithHelp(cfg.Back+"/esc", "back"),
		),
		New: key.NewBinding(
			key.WithKeys(cfg.New),
			key.WithHelp(cfg.New, "new"),
		),
		Select: key.NewBinding(
			key.WithKeys(cfg.Select),
			key.WithHelp(cfg.Select, "select"),
		),
		Upvote: key.NewBinding(
			key.WithKeys(cfg.Upvote),
			key.WithHelp(cfg.Upvote, "upvote"),
		),
		Downvote: key.NewBinding(
			key.WithKeys(cfg.Downvote),
			key.WithHelp(cfg.Downvote, "downvote"),
		),
		Refresh: key.NewBinding(
			key.WithKeys(cfg.Refresh),
			key.WithHelp(cfg.Refresh, "refresh"),
		),
		Login: key.NewBinding(
			key.WithKeys(cfg.Login),
			key.WithHelp(cfg.Login, "login view"),
		),
		Feed: key.NewBinding(
			key.WithKeys(cfg.Feed),
			key.WithHelp(cfg.Feed, "feed view"),
		),
		Communities: key.NewBinding(
			key.WithKeys(cfg.Communities),
			key.WithHelp(cfg.Communities, "communities view"),
		),
		Palette: key.NewBinding(
			key.WithKeys(cfg.Palette),
			key.WithHelp(cfg.Palette, "command palette"),
		),
	}
}

func DefaultKeyMap() KeyMap {
	return NewKeyMap(config.DefaultConfig().Keys)
}
