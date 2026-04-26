package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto/cli/internal/types"
)

type PostDetailModel struct {
	Post     types.Post
	Comments []types.Comment
	Viewport viewport.Model
	Ready    bool
	Width    int
	Height   int
}

func NewPostDetailModel() PostDetailModel {
	return PostDetailModel{
		Viewport: viewport.New(0, 0),
	}
}

func (m PostDetailModel) Init() tea.Cmd {
	return nil
}

func (m PostDetailModel) Update(msg tea.Msg) (PostDetailModel, tea.Cmd) {
	var cmd tea.Cmd
	m.Viewport, cmd = m.Viewport.Update(msg)
	return m, cmd
}

func (m PostDetailModel) View() string {
	if !m.Ready {
		return "Loading post..."
	}
	return m.Viewport.View()
}

func (m *PostDetailModel) SetContent(post types.Post, comments []types.Comment) {
	m.Post = post
	m.Comments = comments
	m.Ready = true
	m.render()
}

func (m *PostDetailModel) SetSize(width, height int) {
	m.Width = width
	m.Height = height
	m.Viewport.Width = width
	m.Viewport.Height = height
	if m.Ready {
		m.render()
	}
}

func (m *PostDetailModel) render() {
	var s strings.Builder

	// Render Post Header
	s.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Render(m.Post.Title) + "\n")
	
	s.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(fmt.Sprintf("u/%s in c/%s • ♥ %d", m.Post.AuthorName, m.Post.CommunityName, m.Post.Score)) + "\n\n")

	// Render Post Content with Glamour
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(m.Width-4),
	)
	out, _ := renderer.Render(m.Post.Content)
	s.WriteString(out + "\n")

	s.WriteString(lipgloss.NewStyle().Bold(true).Render("Comments:") + "\n\n")

	// Simple comment rendering
	for _, c := range m.Comments {
		s.WriteString(m.renderComment(c, 0))
	}

	m.Viewport.SetContent(s.String())
}

func (m PostDetailModel) renderComment(c types.Comment, depth int) string {
	indent := strings.Repeat("  ", depth)
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%s%s %s\n", 
		indent, 
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render(c.AuthorName),
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(fmt.Sprintf("♥ %d", c.Score))))
	
	// Wrap comment content
	contentStyle := lipgloss.NewStyle().PaddingLeft(depth * 2 + 2)
	s.WriteString(contentStyle.Render(c.Content) + "\n\n")

	for _, child := range c.Children {
		s.WriteString(m.renderComment(child, depth+1))
	}
	return s.String()
}
