package views

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingsModel struct {
	Avatar  textinput.Model
	Focused int // 0: Avatar, 1: Save
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
	s += lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render("Profile Settings") + "\n\n"

	s += "Avatar URL:\n" + m.Avatar.View() + "\n\n"

	submitBtn := "[ Save ]"
	if m.Focused == 1 {
		submitBtn = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("[ Save ]")
	}
	s += submitBtn + "\n"

	s += "\n\nTip: Use :delete-account to permanently delete your account."

	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}
