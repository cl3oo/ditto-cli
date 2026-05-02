package views

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type RegisterModel struct {
	Username textinput.Model
	Email    textinput.Model
	Password textinput.Model
	Confirm  textinput.Model
	Focused  int // 0: username, 1: email, 2: password, 3: confirm, 4: submit, 5: cancel
	Error    string
	Theme    theme.Theme
}

func NewRegisterModel() RegisterModel {
	u := textinput.New()
	u.Placeholder = "Username"
	u.Focus()

	e := textinput.New()
	e.Placeholder = "Email"

	p := textinput.New()
	p.Placeholder = "Password"
	p.EchoMode = textinput.EchoPassword
	p.EchoCharacter = '•'

	c := textinput.New()
	c.Placeholder = "Confirm Password"
	c.EchoMode = textinput.EchoPassword
	c.EchoCharacter = '•'

	return RegisterModel{
		Username: u,
		Email:    e,
		Password: p,
		Confirm:  c,
		Focused:  0,
	}
}

func (m RegisterModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *RegisterModel) SetTheme(t theme.Theme) {
	m.Theme = t
}

func (m RegisterModel) Update(msg tea.Msg) (RegisterModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			m.Focused = (m.Focused + 1) % 6
			m.Username.Blur()
			m.Email.Blur()
			m.Password.Blur()
			m.Confirm.Blur()
			switch m.Focused {
			case 0:
				m.Username.Focus()
			case 1:
				m.Email.Focus()
			case 2:
				m.Password.Focus()
			case 3:
				m.Confirm.Focus()
			}
		}
	}

	switch m.Focused {
	case 0:
		m.Username, cmd = m.Username.Update(msg)
	case 1:
		m.Email, cmd = m.Email.Update(msg)
	case 2:
		m.Password, cmd = m.Password.Update(msg)
	case 3:
		m.Confirm, cmd = m.Confirm.Update(msg)
	}

	return m, cmd
}

func (m RegisterModel) View() string {
	var s string

	s += m.Theme.AccentText.Bold(true).Render("Register a New Ditto Account") + "\n\n"
	s += m.Username.View() + "\n"
	s += m.Email.View() + "\n"
	s += m.Password.View() + "\n"
	s += m.Confirm.View() + "\n\n"

	submitBtn := m.Theme.RenderPrimaryButton("[ Register ]", m.Focused == 4)
	cancelBtn := m.Theme.RenderSecondaryButton("[ Cancel ]", m.Focused == 5)

	s += submitBtn + "  " + cancelBtn + "\n"

	if m.Error != "" {
		s += "\n" + m.Theme.Error.Render(m.Error)
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}
