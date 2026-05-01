package views

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto/cli/internal/types"
)

type PostItem struct {
	types.Post
}

func (i PostItem) Title() string       { return i.Post.Title }
func (i PostItem) Description() string {
	return fmt.Sprintf("u/%s in c/%s • ↑↓ %d • 💬 %d", 
		i.Post.Author.Username, 
		i.Post.Community.Name, 
		i.Post.Scores.VoteScore,
		i.Post.Scores.CommentCount)
}
func (i PostItem) FilterValue() string { return i.Post.Title + " " + i.Post.Author.Username }

type itemDelegate struct {
	Theme lipgloss.Style // Selected style
}

func (d itemDelegate) Height() int                               { return 2 }
func (d itemDelegate) Spacing() int                              { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(PostItem)
	if !ok {
		return
	}

	title := lipgloss.NewStyle().Bold(true).Render(i.Title())
	desc := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(i.Description())

	if index == m.Index() {
		title = d.Theme.Render(i.Title())
	}

	fmt.Fprintf(w, "  %s\n  %s", title, desc)
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
