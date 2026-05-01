package views

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/types"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type PostItem struct {
	types.Post
}

func (i PostItem) Title() string { return i.Post.Title }
func (i PostItem) Description() string {
	return fmt.Sprintf("u/%s in c/%s • ↑↓ %d • 💬 %d",
		i.Author.Username,
		i.Community.Name,
		i.Scores.VoteScore,
		i.Scores.CommentCount)
}
func (i PostItem) FilterValue() string { return i.Post.Title + " " + i.Author.Username }

func RelativeTime(t time.Time) string {
	d := time.Since(t)
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	return t.Format("Jan 02")
}

type itemDelegate struct {
	Theme theme.Theme
}

func (d itemDelegate) Height() int                               { return 5 }
func (d itemDelegate) Spacing() int                              { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(PostItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	cardStyle := d.Theme.Post.
		Width(max(20, m.Width()-4))

	if isSelected {
		cardStyle = cardStyle.BorderForeground(d.Theme.Accent)
	}

	header := d.Theme.TextSubtle.Render(
		fmt.Sprintf("u/%s in c/%s • %s", i.Author.Username, i.Community.Name, RelativeTime(i.CreatedAt)),
	)

	titleStyle := d.Theme.Text.Bold(true)
	if isSelected {
		titleStyle = d.Theme.AccentText.Bold(true)
	}
	title := titleStyle.Render(i.Title())

	stats := d.Theme.TextSubtle.Render(
		fmt.Sprintf("↑↓ %d • 💬 %d • 💎 %d",
			i.Scores.VoteScore,
			i.Scores.CommentCount,
			i.Scores.AwardCount),
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		title,
		"",
		stats,
	)

	_, _ = fmt.Fprint(w, cardStyle.Render(content))
}

type FeedModel struct {
	List   list.Model
	Loaded bool
	Theme  theme.Theme
}

func NewFeedModel() FeedModel {
	l := list.New([]list.Item{}, itemDelegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return FeedModel{
		List: l,
	}
}

func (m *FeedModel) SetTheme(theme theme.Theme) {
	m.Theme = theme
	m.List.SetDelegate(itemDelegate{Theme: theme})
}
func (m FeedModel) Init() tea.Cmd {
	return nil
}

func (m FeedModel) Update(msg tea.Msg) (FeedModel, tea.Cmd) {
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m FeedModel) View() string {
	return m.List.View()
}

func (m *FeedModel) SetPosts(posts []types.Post) {
	items := make([]list.Item, len(posts))
	for i, post := range posts {
		items[i] = PostItem{Post: post}
	}
	m.List.SetItems(items)
	m.Loaded = true
}

func (m *FeedModel) SetSize(width, height int) {
	m.List.SetSize(width, height)
}
