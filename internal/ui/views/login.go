package views

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	Width        int
}

func NewLoginModel() LoginModel {
	u := textinput.New()
	u.Placeholder = "Username"
	u.Focus()

	p := textinput.New()
	p.Placeholder = "Password"
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'
	u.Width = fitInputWidth(0, 40)
	p.Width = fitInputWidth(0, 40)

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

func (m *LoginModel) SetSize(width int) {
	m.Width = width
	inputWidth := fitInputWidth(width, 40)
	m.Username.Width = inputWidth
	m.Password.Width = inputWidth
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
	title := "Login to Ditto"
	subtitle := "Use your Ditto account to unlock the TUI."
	sections := []formSection{}
	var feedback string

	if m.LoggedIn {
		title = "You are in"
		subtitle = "The session token is already saved locally."
	}

	if !m.LoggedIn {
		sections = append(sections,
			formSection{Label: "Username", Content: m.Username.View()},
			formSection{Label: "Password", Content: m.Password.View()},
		)
	} else {
		displayToken := m.SuccessToken
		if len(displayToken) > 8 {
			displayToken = displayToken[:8] + "..."
		}
		feedback = renderFormFeedback(m.Theme, "success", "Session saved", "Welcome back. Token preview: "+displayToken)
	}

	submitLabel := "[ Log In ]"
	if m.LoggedIn {
		submitLabel = "[ Continue ]"
	}

	submitBtn := m.Theme.RenderPrimaryButton(submitLabel, m.Focused == 2)
	registerBtn := m.Theme.RenderSecondaryButton("[ Create Account ]", m.Focused == 3)

	if m.Error != "" {
		feedback = renderFormFeedback(m.Theme, "error", "Could not log in", fmt.Sprintf("Error: %s", m.Error))
	}

	return renderFormShell(m.Theme, m.Width, title, subtitle, sections, []string{submitBtn, registerBtn}, feedback)
}
