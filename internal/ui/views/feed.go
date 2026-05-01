package views

import (
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/types"
)

type PostItem struct {
	types.Post
}

func (i PostItem) Title() string       { return i.Post.Title }
func (i PostItem) Description() string { return i.Post.Content }
func (i PostItem) FilterValue() string { return i.Post.Title + " " + i.Post.Author.Username }

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
	Theme lipgloss.Style // Selected style
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
	
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(0, 1).
		Width(m.Width() - 4)

	if isSelected {
		cardStyle = cardStyle.BorderForeground(lipgloss.Color("170"))
	}

	// Header: u/user in c/community • time
	header := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
		fmt.Sprintf("u/%s in c/%s • %s", i.Author.Username, i.Community.Name, RelativeTime(i.CreatedAt)),
	)

	// Title: Bold and colorful
	titleStyle := lipgloss.NewStyle().Bold(true)
	if isSelected {
		titleStyle = titleStyle.Foreground(lipgloss.Color("170"))
	}
	title := titleStyle.Render(i.Title())

	// Footer: stats
	stats := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
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

	fmt.Fprint(w, cardStyle.Render(content))
}

type FeedModel struct {
	List   list.Model
	Loaded bool
	Theme  lipgloss.Style
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

func (m *FeedModel) SetTheme(selected lipgloss.Style) {
	m.Theme = selected
	m.List.SetDelegate(itemDelegate{Theme: selected})
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
