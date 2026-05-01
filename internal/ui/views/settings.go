package views

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type SettingsModel struct {
	Avatar  textinput.Model
	Focused int // 0: Avatar, 1: Save
	Theme   theme.Theme
}

func NewSettingsModel() SettingsModel {
	a := textinput.New()
	a.Placeholder = "Avatar URL"
	a.Focus()

	return SettingsModel{
		Avatar:  a,
		Focused: 0,
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m *SettingsModel) SetTheme(t theme.Theme) {
	m.Theme = t
}

func (m SettingsModel) Update(msg tea.Msg) (SettingsModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			if m.Focused == 0 {
				m.Focused = 1
				m.Avatar.Blur()
			} else {
				m.Focused = 0
				m.Avatar.Focus()
			}
		}
	}

	m.Avatar, cmd = m.Avatar.Update(msg)
	return m, cmd
}

func (m SettingsModel) View() string {
	var s string
	s += m.Theme.AccentText.Bold(true).Render("Profile Settings") + "\n\n"

	s += m.Theme.Text.Render("Avatar URL:") + "\n" + m.Avatar.View() + "\n\n"

	submitBtn := "[ Save ]"
	if m.Focused == 1 {
		submitBtn = m.Theme.Selected.Render("[ Save ]")
	}
	s += submitBtn + "\n"

	s += "\n\n" + m.Theme.TextSubtle.Render("Tip: Use :delete-account to permanently delete your account.")

	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}
