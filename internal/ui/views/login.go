package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type LoginModel struct {
	Username     textinput.Model
	Password     textinput.Model
	Focused      int // 0 for username, 1 for password, 2 for submit, 3 for register
	Error        string
	SuccessToken string
	LoggedIn     bool
	Theme        theme.Theme
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

func (m *LoginModel) SetTheme(t theme.Theme) {
	m.Theme = t
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
	s += m.Theme.AccentText.Bold(true).Render(title) + "\n\n"

	if !m.LoggedIn {
		s += m.Username.View() + "\n"
		s += m.Password.View() + "\n\n"
	} else {
		displayToken := m.SuccessToken
		if len(displayToken) > 8 {
			displayToken = displayToken[:8] + "..."
		}
		s += m.Theme.Success.Render("Welcome back! Your session token was saved.") + "\n"
		s += m.Theme.TextSubtle.Render(displayToken) + "\n\n"
	}

	submitLabel := "[ Submit ]"
	if m.LoggedIn {
		submitLabel = "[ Save & Continue ]"
	}

	submitBtn := m.Theme.RenderPrimaryButton(submitLabel, m.Focused == 2)
	registerBtn := m.Theme.RenderSecondaryButton("[ Register New Account ]", m.Focused == 3)

	s += submitBtn + "  " + registerBtn + "\n"

	if m.Error != "" {
		cuteMsg := "Oopsie! Something went wrong... (´･ω･`)"
		s += "\n" + m.Theme.TextSubtle.Render(cuteMsg) + "\n"
		s += m.Theme.Error.Render(fmt.Sprintf("Error: %s", m.Error))
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}
