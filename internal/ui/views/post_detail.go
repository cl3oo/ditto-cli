package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto/cli/internal/types"
)

type PostDetailTheme struct {
	Accent   lipgloss.Color
	Selected lipgloss.Style
	Markdown string
}

type PostDetailModel struct {
	Post             types.Post
	Comments         []types.Comment
	Viewport         viewport.Model
	CommentInput     textinput.Model
	ShowCommentInput bool
	Ready            bool
	Width            int
	Height           int
	Theme            PostDetailTheme
}

func NewPostDetailModel() PostDetailModel {
	ti := textinput.New()
	ti.Placeholder = "Write a comment..."
	ti.CharLimit = 1000
	ti.Width = 50

	return PostDetailModel{
		Viewport:     viewport.New(0, 0),
		CommentInput: ti,
	}
}

func (m *PostDetailModel) SetTheme(accent lipgloss.Color, selected lipgloss.Style, markdown string) {
	m.Theme = PostDetailTheme{Accent: accent, Selected: selected, Markdown: markdown}
}

func (m PostDetailModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PostDetailModel) Update(msg tea.Msg) (PostDetailModel, tea.Cmd) {
	var cmd, tiCmd, vpCmd tea.Cmd

	if m.ShowCommentInput {
		m.CommentInput, tiCmd = m.CommentInput.Update(msg)
	}
	m.Viewport, vpCmd = m.Viewport.Update(msg)

	cmd = tea.Batch(tiCmd, vpCmd)
	return m, cmd
}

func (m PostDetailModel) View() string {
	if !m.Ready {
		return "  Loading post..."
	}
	
	v := m.Viewport.View()
	if m.ShowCommentInput {
		return lipgloss.JoinVertical(lipgloss.Left,
			v,
			"\n"+lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), true, false, false, false).
				BorderForeground(lipgloss.Color("241")).
				Padding(1, 1).
				Render(m.CommentInput.View()),
		)
	}
	return v
}

func (m *PostDetailModel) SetContent(post types.Post, comments []types.Comment) {
	m.Post = post
	m.Comments = comments
	m.Ready = true
	m.render()
}

func (m *PostDetailModel) SetShowCommentInput(show bool) {
	m.ShowCommentInput = show
	if show {
		m.CommentInput.Focus()
	} else {
		m.CommentInput.Blur()
	}
	m.SetSize(m.Width, m.Height)
}

func (m *PostDetailModel) SetSize(width, height int) {
	m.Width = width
	m.Height = height
	m.Viewport.Width = width
	if m.ShowCommentInput {
		m.Viewport.Height = height - 5 // Reserved for comment input
	} else {
		m.Viewport.Height = height
	}
	m.CommentInput.Width = width - 4
	if m.Ready {
		m.render()
	}
}

func (m *PostDetailModel) render() {
	var s strings.Builder

	// Render Post Header
	s.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Accent).
		Render(fmt.Sprintf("%s (p/%s)", m.Post.Title, m.Post.ID)) + "\n")
	
	s.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(fmt.Sprintf("u/%s (u/%s) in c/%s (c/%s) • ♥ %d", 
			m.Post.Author.Username, m.Post.Author.ID,
			m.Post.Community.Name, m.Post.Community.ID,
			m.Post.Scores.VoteScore)) + "\n\n")

	// Render Post Content with Glamour
	var renderer *glamour.TermRenderer
	if m.Theme.Markdown != "" {
		renderer, _ = glamour.NewTermRenderer(
			glamour.WithStylesFromJSONBytes([]byte(m.Theme.Markdown)),
			glamour.WithWordWrap(m.Width-4),
		)
	} else {
		renderer, _ = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(m.Width-4),
		)
	}
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
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render(c.Author.Username),
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(fmt.Sprintf("♥ %d", c.Scores.VoteScore))))
	
	// Wrap comment content
	contentStyle := lipgloss.NewStyle().PaddingLeft(depth * 2 + 2)
	s.WriteString(contentStyle.Render(c.Content) + "\n\n")

	for _, child := range c.Children {
		s.WriteString(m.renderComment(child, depth+1))
	}
	return s.String()
}
