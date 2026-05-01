package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoginModel struct {
	Username     textinput.Model
	Password     textinput.Model
	Focused      int // 0 for username, 1 for password, 2 for submit, 3 for register
	Error        string
	SuccessToken string
	LoggedIn     bool
}

func NewLoginModel() LoginModel {
	u := textinput.New()
	u.Placeholder = "Username"
	u.Focus()

	p := textinput.New()
	p.Placeholder = "Password"
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'

	return LoginModel{
		Username: u,
		Password: p,
		Focused:  0,
	}
}

func (m LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LoginModel) Update(msg tea.Msg) (LoginModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			m.Focused = (m.Focused + 1) % 4
			switch m.Focused {
			case 0:
				m.Username.Focus()
				m.Password.Blur()
			case 1:
				m.Username.Blur()
				m.Password.Focus()
			default:
				m.Username.Blur()
				m.Password.Blur()
			}
		}
	}

	switch m.Focused {
	case 0:
		m.Username, cmd = m.Username.Update(msg)
	case 1:
		m.Password, cmd = m.Password.Update(msg)
	}

	return m, cmd
}

func (m LoginModel) View() string {
	var s string

	title := "Login to Ditto"
	if m.LoggedIn {
		title = "Successfully Authenticated! ✨"
	}
	s += lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render(title) + "\n\n"

	if !m.LoggedIn {
		s += m.Username.View() + "\n"
		s += m.Password.View() + "\n\n"
	} else {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Render("Welcome back! Here is your token:") + "\n"
		s += lipgloss.NewStyle().Italic(true).Faint(true).Render(m.SuccessToken) + "\n\n"
	}

	submitLabel := "[ Submit ]"
	if m.LoggedIn {
		submitLabel = "[ Save & Continue ]"
	}

	submitBtn := submitLabel
	if m.Focused == 2 {
		submitBtn = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render(submitLabel)
	}

	registerBtn := "[ Register New Account ]"
	if m.Focused == 3 {
		registerBtn = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("[ Register New Account ]")
	}

	s += submitBtn + "  " + registerBtn + "\n"

	if m.Error != "" {
		cuteMsg := "Oopsie! Something went wrong... (´･ω･`)"
		s += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("211")).Render(cuteMsg) + "\n"
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("Error: %s", m.Error))
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}
