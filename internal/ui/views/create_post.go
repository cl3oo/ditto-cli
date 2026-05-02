package views

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type CreatePostModel struct {
	Title       textinput.Model
	CommunityID textinput.Model
	Content     textarea.Model
	Focused     int // 0: Title, 1: CommunityID, 2: Content, 3: Submit
	Theme       theme.Theme
	Width       int
}

func NewCreatePostModel() CreatePostModel {
	t := textinput.New()
	t.Placeholder = "Title"
	t.Focus()

	cID := textinput.New()
	cID.Placeholder = "Community Name"

	cont := textarea.New()
	cont.Placeholder = "Content (Markdown supported)..."
	fieldWidth := fitInputWidth(0, 60)
	t.Width = fieldWidth
	cID.Width = fieldWidth
	cont.SetWidth(fieldWidth)
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

func (m *CreatePostModel) SetSize(width int) {
	m.Width = width
	fieldWidth := fitInputWidth(width, 60)
	m.Title.Width = fieldWidth
	m.CommunityID.Width = fieldWidth
	m.Content.SetWidth(fieldWidth)
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
	submitBtn := m.Theme.RenderPrimaryButton("[ Submit ]", m.Focused == 3)
	sections := []formSection{
		{Label: "Title", Content: m.Title.View()},
		{Label: "Community", Content: m.CommunityID.View()},
		{Label: "Content", Content: m.Content.View()},
	}
	return renderFormShell(m.Theme, m.Width, "Create new post", "Draft a post without leaving the terminal.", sections, []string{submitBtn}, "")
}
