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
	FlattenedComments []commentWithDepth
	SelectedIdx      int // -1 for post, 0+ for comments
	
	Viewport         viewport.Model
	CommentInput     textinput.Model
	ShowCommentInput bool
	Ready            bool
	Width            int
	Height           int
	Theme            PostDetailTheme
}

type commentWithDepth struct {
	types.Comment
	Depth int
}

func NewPostDetailModel() PostDetailModel {
	ti := textinput.New()
	ti.Placeholder = "Write a comment..."
	ti.CharLimit = 1000
	ti.Width = 50

	return PostDetailModel{
		Viewport:     viewport.New(0, 0),
		CommentInput: ti,
		SelectedIdx:  -1,
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
		return m, tiCmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			if m.SelectedIdx < len(m.FlattenedComments)-1 {
				m.SelectedIdx++
				m.render()
			}
		case "k":
			if m.SelectedIdx > -1 {
				m.SelectedIdx--
				m.render()
			}
		}
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
	m.FlattenedComments = []commentWithDepth{}
	m.flattenComments(comments, 0)
	m.Ready = true
	m.SelectedIdx = -1
	m.render()
}

func (m *PostDetailModel) flattenComments(comments []types.Comment, depth int) {
	for _, c := range comments {
		m.FlattenedComments = append(m.FlattenedComments, commentWithDepth{Comment: c, Depth: depth})
		m.flattenComments(c.Children, depth+1)
	}
}

func (m *PostDetailModel) SetShowCommentInput(show bool) {
	if show && m.Post.Locked {
		return
	}
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
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.Theme.Accent)
	
	if m.SelectedIdx == -1 {
		headerStyle = headerStyle.Background(lipgloss.Color("235"))
	}
	
	headerText := fmt.Sprintf("%s (p/%s)", m.Post.Title, m.Post.ID)
	if m.Post.Locked {
		headerText += " [LOCKED]"
	}
	s.WriteString(headerStyle.Render(headerText) + "\n")
	
	s.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(fmt.Sprintf("u/%s (u/%s) in c/%s (c/%s) • ↑↓ %d • 💬 %d", 
			m.Post.Author.Username, m.Post.Author.ID,
			m.Post.Community.Name, m.Post.Community.ID,
			m.Post.Scores.VoteScore,
			m.Post.Scores.CommentCount)) + "\n\n")

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

	// Render flattened comments with selection
	for i, c := range m.FlattenedComments {
		s.WriteString(m.renderCommentItem(c, i == m.SelectedIdx))
	}

	m.Viewport.SetContent(s.String())
}

func (m PostDetailModel) renderCommentItem(c commentWithDepth, selected bool) string {
	indent := strings.Repeat("│ ", c.Depth)
	if c.Depth > 0 {
		indent = strings.Repeat("│ ", c.Depth-1) + "├─"
	}
	
	var s strings.Builder
	authorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
	scoreStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	contentStyle := lipgloss.NewStyle().PaddingLeft(c.Depth * 2)

	if selected {
		authorStyle = authorStyle.Background(lipgloss.Color("235"))
		contentStyle = contentStyle.Background(lipgloss.Color("235"))
	}

	s.WriteString(fmt.Sprintf("%s %s %s\n", 
		indent, 
		authorStyle.Render("u/"+c.Author.Username),
		scoreStyle.Render(fmt.Sprintf("↑↓ %d", c.Scores.VoteScore))))
	
	// Wrap comment content
	s.WriteString(contentStyle.Render(c.Content) + "\n\n")

	return s.String()
}
