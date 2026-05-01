package views

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RegisterModel struct {
	Username textinput.Model
	Email    textinput.Model
	Password textinput.Model
	Confirm  textinput.Model
	Focused  int // 0: username, 1: email, 2: password, 3: confirm, 4: submit, 5: cancel
	Error    string
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
			if m.Focused == 0 {
				m.Username.Focus()
			} else if m.Focused == 1 {
				m.Email.Focus()
			} else if m.Focused == 2 {
				m.Password.Focus()
			} else if m.Focused == 3 {
				m.Confirm.Focus()
			}
		}
	}

	if m.Focused == 0 {
		m.Username, cmd = m.Username.Update(msg)
	} else if m.Focused == 1 {
		m.Email, cmd = m.Email.Update(msg)
	} else if m.Focused == 2 {
		m.Password, cmd = m.Password.Update(msg)
	} else if m.Focused == 3 {
		m.Confirm, cmd = m.Confirm.Update(msg)
	}

	return m, cmd
}

func (m RegisterModel) View() string {
	var s string

	s += lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render("Register a New Ditto Account") + "\n\n"
	s += m.Username.View() + "\n"
	s += m.Email.View() + "\n"
	s += m.Password.View() + "\n"
	s += m.Confirm.View() + "\n\n"

	submitBtn := "[ Register ]"
	if m.Focused == 4 {
		submitBtn = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true).Render("[ Register ]")
	}

	cancelBtn := "[ Cancel ]"
	if m.Focused == 5 {
		cancelBtn = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Bold(true).Render("[ Cancel ]")
	}
	
	s += submitBtn + "  " + cancelBtn + "\n"

	if m.Error != "" {
		s += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(m.Error)
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}
