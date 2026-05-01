package views

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type CreatePostModel struct {
	Title       textinput.Model
	CommunityID textinput.Model
	Content     textarea.Model
	Focused     int // 0: Title, 1: CommunityID, 2: Content, 3: Submit
	Theme       theme.Theme
}

func NewCreatePostModel() CreatePostModel {
	t := textinput.New()
	t.Placeholder = "Title"
	t.Focus()

	cID := textinput.New()
	cID.Placeholder = "Community Name"

	cont := textarea.New()
	cont.Placeholder = "Content (Markdown supported)..."
	cont.SetWidth(60)
	cont.SetHeight(10)

	return CreatePostModel{
		Title:       t,
		CommunityID: cID,
		Content:     cont,
		Focused:     0,
	}
}

func (m CreatePostModel) Init() tea.Cmd {
	return nil
}

func (m *CreatePostModel) SetTheme(t theme.Theme) {
	m.Theme = t
}

func (m CreatePostModel) Update(msg tea.Msg) (CreatePostModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			if msg.String() == "tab" {
				m.Focused = (m.Focused + 1) % 4
			} else {
				m.Focused = (m.Focused - 1 + 4) % 4
			}

			m.Title.Blur()
			m.CommunityID.Blur()
			m.Content.Blur()

			switch m.Focused {
			case 0:
				m.Title.Focus()
			case 1:
				m.CommunityID.Focus()
			case 2:
				m.Content.Focus()
			}
		}
	}

	m.Title, cmd = m.Title.Update(msg)
	cmds = append(cmds, cmd)
	m.CommunityID, cmd = m.CommunityID.Update(msg)
	cmds = append(cmds, cmd)
	m.Content, cmd = m.Content.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m CreatePostModel) View() string {
	var s string
	s += m.Theme.AccentText.Bold(true).Render("Create New Post") + "\n\n"

	s += m.Theme.Text.Render("Title:") + "\n" + m.Title.View() + "\n\n"
	s += m.Theme.Text.Render("Community:") + "\n" + m.CommunityID.View() + "\n\n"
	s += m.Theme.Text.Render("Content:") + "\n" + m.Content.View() + "\n\n"

	submitBtn := "[ Submit ]"
	if m.Focused == 3 {
		submitBtn = m.Theme.Selected.Render("[ Submit ]")
	}
	s += submitBtn + "\n"

	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}
