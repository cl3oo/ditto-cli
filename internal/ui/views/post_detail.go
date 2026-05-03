package views

import (
	"fmt"
	"strings"
	"time"

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
	pendingG         bool
}

type commentWithDepth struct {
	types.Comment
	Depth        int
	AncestorLast []bool
	IsLast       bool
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
			m.pendingG = false
		case "k":
			if m.SelectedIdx > -1 {
				m.SelectedIdx--
				m.render()
			}
			m.pendingG = false
		case "g":
			if m.pendingG {
				m.SelectedIdx = -1
				m.render()
				m.pendingG = false
				return m, nil
			}
			m.pendingG = true
			return m, nil
		case "G":
			m.SelectedIdx = len(m.FlattenedComments) - 1
			m.render()
			m.pendingG = false
		case "H":
			// Jump to top
			m.SelectedIdx = -1
			m.render()
			m.pendingG = false
		case "L":
			// Jump to bottom
			m.SelectedIdx = len(m.FlattenedComments) - 1
			m.render()
			m.pendingG = false
		case "M":
			// Jump to middle
			m.SelectedIdx = (len(m.FlattenedComments) - 1) / 2
			m.render()
			m.pendingG = false
		case "c":
			m.PrepareReplyInput(false)
			m.pendingG = false
			return m, nil
		case "enter":
			if m.SelectedIdx >= 0 {
				m.PrepareReplyInput(true)
				m.pendingG = false
				return m, nil
			}
		default:
			m.pendingG = false
		}
	}

	m.Viewport, vpCmd = m.Viewport.Update(msg)

	cmd = tea.Batch(tiCmd, vpCmd)
	return m, cmd
}

func (m PostDetailModel) View() string {
	if !m.Ready {
		return RenderCenteredStateSurface(m.Theme, m.Width, m.Height, "loading", "Loading post", "Fetching the thread and replies.", nil)
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
	m.rebuildFlattenedComments()
	m.Ready = true
	m.SelectedIdx = -1
	m.Page = 1
	m.render()
}

func (m *PostDetailModel) AppendComments(comments []types.Comment) {
	m.Comments = mergeCommentForest(m.Comments, comments)
	m.rebuildFlattenedComments()
	m.Page++
	m.render()
}

func (m *PostDetailModel) rebuildFlattenedComments() {
	m.FlattenedComments = m.FlattenedComments[:0]
	m.flattenComments(m.Comments, 0, nil)
}

func (m *PostDetailModel) flattenComments(comments []types.Comment, depth int, ancestorLast []bool) {
	for i, c := range comments {
		isLast := i == len(comments)-1
		prefix := append([]bool(nil), ancestorLast...)
		m.FlattenedComments = append(m.FlattenedComments, commentWithDepth{
			Comment:      c,
			Depth:        depth,
			AncestorLast: prefix,
			IsLast:       isLast,
		})
		m.flattenComments(c.Children, depth+1, append(prefix, isLast))
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
		fmt.Sprintf("Posted %s", formatDetailTimestamp(m.Post.CreatedAt)),
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
	var s strings.Builder
	authorStyle := m.Theme.Text.Bold(true)
	scoreStyle := m.Theme.TextSubtle

	if selected {
		authorStyle = m.Theme.AccentText.Bold(true)
	}

	headerPrefix, contentPrefix := commentTreePrefixes(c)

	_, _ = fmt.Fprintf(&s, "%s %s • %s • %s\n",
		headerPrefix,
		authorStyle.Render("u/"+c.Author.Username),
		scoreStyle.Render(fmt.Sprintf("↑↓ %d", c.Scores.VoteScore)),
		scoreStyle.Render(RelativeTime(c.CreatedAt)))

	detailParts := []string{formatDetailTimestamp(c.CreatedAt)}
	if replies := len(c.Children); replies > 0 {
		label := "replies"
		if replies == 1 {
			label = "reply"
		}
		detailParts = append(detailParts, fmt.Sprintf("%d %s", replies, label))
	}
	_, _ = fmt.Fprintf(&s, "%s %s\n", contentPrefix, scoreStyle.Render(strings.Join(detailParts, " • ")))

	contentWidth := max(12, m.Width-lipgloss.Width(contentPrefix)-4)
	wrapped := wrapCommentContent(c.Content, contentWidth)
	for i, line := range wrapped {
		if i > 0 {
			s.WriteByte('\n')
		}
		s.WriteString(contentPrefix)
		s.WriteString(line)
	}

	return withSelectionIndicator(s.String(), selected, m.Theme) + "\n\n"
}

func formatDetailTimestamp(t time.Time) string {
	if t.IsZero() {
		return "unknown time"
	}
	return t.Local().Format("2006-01-02 15:04")
}

func commentTreePrefixes(c commentWithDepth) (string, string) {
	if c.Depth == 0 {
		return "•", "  "
	}

	var shared strings.Builder
	for _, ancestorIsLast := range c.AncestorLast[:len(c.AncestorLast)-1] {
		if ancestorIsLast {
			shared.WriteString("  ")
		} else {
			shared.WriteString("│ ")
		}
	}

	branch := "├─"
	stem := "│ "
	if c.IsLast {
		branch = "└─"
		stem = "  "
	}

	return shared.String() + branch, shared.String() + stem
}

func wrapCommentContent(content string, width int) []string {
	content = strings.TrimSpace(content)
	if content == "" {
		return []string{""}
	}

	paragraphs := strings.Split(content, "\n")
	lines := make([]string, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			if len(lines) == 0 || lines[len(lines)-1] != "" {
				lines = append(lines, "")
			}
			continue
		}

		words := strings.Fields(paragraph)
		line := words[0]
		for _, word := range words[1:] {
			candidate := line + " " + word
			if lipgloss.Width(candidate) <= width {
				line = candidate
				continue
			}
			lines = append(lines, line)
			line = word
		}
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return []string{""}
	}

	return lines
}

func mergeCommentForest(existing, incoming []types.Comment) []types.Comment {
	merged := append([]types.Comment(nil), existing...)
	indexByID := make(map[string]int, len(merged))
	for i, comment := range merged {
		indexByID[comment.ID] = i
	}

	for _, comment := range incoming {
		if idx, ok := indexByID[comment.ID]; ok {
			merged[idx] = mergeCommentNode(merged[idx], comment)
			continue
		}
		indexByID[comment.ID] = len(merged)
		merged = append(merged, comment)
	}

	return merged
}

func mergeCommentNode(existing, incoming types.Comment) types.Comment {
	merged := existing
	merged.Content = incoming.Content
	merged.AuthorID = incoming.AuthorID
	merged.Author = incoming.Author
	merged.Scores = incoming.Scores
	merged.Voted = incoming.Voted
	merged.VoteDirection = incoming.VoteDirection
	merged.CreatedAt = incoming.CreatedAt
	merged.Children = mergeCommentForest(existing.Children, incoming.Children)
	return merged
}
