package views

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type SettingsModel struct {
	Avatar  textinput.Model
	Focused int // 0: Avatar, 1: Save
	Theme   theme.Theme
	Width   int
}

func NewSettingsModel() SettingsModel {
	a := textinput.New()
	a.Placeholder = "Avatar URL"
	a.Focus()
	a.Width = fitInputWidth(0, 40)

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

func (m *SettingsModel) SetSize(width int) {
	m.Width = width
	m.Avatar.Width = fitInputWidth(width, 40)
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
	submitBtn := m.Theme.RenderPrimaryButton("[ Save ]", m.Focused == 1)
	sections := []formSection{{Label: "Avatar URL", Content: m.Avatar.View()}}
	feedback := renderFormFeedback(m.Theme, "error", "Danger zone", "Use :delete-account to permanently remove your account.")
	return renderFormShell(m.Theme, m.Width, "Profile settings", "Update the account details that the TUI can edit today.", sections, []string{submitBtn}, feedback)
}
