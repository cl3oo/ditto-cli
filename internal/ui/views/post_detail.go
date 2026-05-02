package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/types"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type PostDetailModel struct {
	Post              types.Post
	Comments          []types.Comment
	FlattenedComments []commentWithDepth
	SelectedIdx       int // -1 for post, 0+ for comments

	Viewport         viewport.Model
	CommentInput     textinput.Model
	ShowCommentInput bool
	Ready            bool
	Page             int
	Width            int
	Height           int
	Theme            theme.Theme
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

func (m *PostDetailModel) SetTheme(t theme.Theme) {
	m.Theme = t
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
		case "c":
			m.PrepareReplyInput(false)
			return m, nil
		case "enter":
			if m.SelectedIdx >= 0 {
				m.PrepareReplyInput(true)
				return m, nil
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

	var actionMenu string
	if !m.ShowCommentInput {
		if m.SelectedIdx == -1 {
			actionMenu = m.Theme.TextSubtle.Render("Actions: (c)omment • (r)eply • (p)rofile • (s)hare • (R)eport • (L)oad more")
		} else {
			author := m.FlattenedComments[m.SelectedIdx].Author.Username
			actionMenu = m.Theme.AccentText.Bold(true).Render(fmt.Sprintf("Comment by u/%s: (enter) reply • (r)eply • (p)rofile • (s)hare • (R)eport", author))
		}
	}

	v := m.Viewport.View()

	var footer string
	if m.ShowCommentInput {
		footer = "\n" + m.Theme.Surface.
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(m.Theme.TextSubtle.GetForeground()).
			Padding(1, 1).
			Render(m.CommentInput.View())
	} else {
		footer = "\n" + lipgloss.NewStyle().
			Padding(0, 2).
			Render(actionMenu)
	}

	return lipgloss.JoinVertical(lipgloss.Left, v, footer)
}

func (m *PostDetailModel) SetContent(post types.Post, comments []types.Comment) {
	m.Post = post
	m.Comments = comments
	m.FlattenedComments = []commentWithDepth{}
	m.flattenComments(comments, 0)
	m.Ready = true
	m.SelectedIdx = -1
	m.Page = 1
	m.render()
}

func (m *PostDetailModel) AppendComments(comments []types.Comment) {
	m.Comments = append(m.Comments, comments...)
	m.flattenComments(comments, 0)
	m.Page++
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

func (m *PostDetailModel) PrepareReplyInput(replyToSelectedComment bool) {
	if replyToSelectedComment && m.SelectedIdx >= 0 && m.SelectedIdx < len(m.FlattenedComments) {
		author := m.FlattenedComments[m.SelectedIdx].Author.Username
		m.CommentInput.Placeholder = "Replying to u/" + author + "..."
	} else {
		m.CommentInput.Placeholder = "Write a comment..."
	}
	m.SetShowCommentInput(true)
}

func (m *PostDetailModel) SetSize(width, height int) {
	m.Width = width
	m.Height = height
	m.Viewport.Width = width
	m.Viewport.Height = height - 3 // Reserved for footer
	m.CommentInput.Width = width - 4
	if m.Ready {
		m.render()
	}
}

func (m *PostDetailModel) render() {
	var s strings.Builder

	headerText := m.Post.Title
	if m.Post.Locked {
		headerText += " [locked]"
	}
	s.WriteString(renderVerticalCard(
		lipgloss.NewStyle(),
		m.SelectedIdx == -1,
		m.Theme,
		fmt.Sprintf("c/%s • u/%s • %s • p/%s", m.Post.Community.Name, m.Post.Author.Username, RelativeTime(m.Post.CreatedAt), m.Post.ID),
		headerText,
		"",
		fmt.Sprintf("↑↓ %d • 💬 %d • 💎 %d", m.Post.Scores.VoteScore, m.Post.Scores.CommentCount, m.Post.Scores.AwardCount),
	) + "\n\n")

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

	s.WriteString(m.Theme.Text.Bold(true).Render("Comments:") + "\n\n")

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
	authorStyle := m.Theme.Text.Bold(true)
	scoreStyle := m.Theme.TextSubtle
	contentStyle := lipgloss.NewStyle().PaddingLeft(c.Depth*2 + 2)

	if selected {
		authorStyle = m.Theme.AccentText.Bold(true)
	}

	_, _ = fmt.Fprintf(&s, "%s %s • %s • %s\n",
		indent,
		authorStyle.Render("u/"+c.Author.Username),
		scoreStyle.Render(fmt.Sprintf("↑↓ %d", c.Scores.VoteScore)),
		scoreStyle.Render(RelativeTime(c.CreatedAt)))

	// Wrap comment content
	s.WriteString(contentStyle.Render(c.Content))

	return withSelectionIndicator(s.String(), selected, m.Theme) + "\n\n"
}
