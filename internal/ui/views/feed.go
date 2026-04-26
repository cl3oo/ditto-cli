package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto/cli/internal/types"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type PostItem struct {
	types.Post
}

func (i PostItem) Title() string       { return i.Post.Title }
func (i PostItem) Description() string { return fmt.Sprintf("u/%s in c/%s • ♥ %d", i.Post.AuthorName, i.Post.CommunityName, i.Post.Score) }
func (i PostItem) FilterValue() string { return i.Post.Title + " " + i.Post.AuthorName }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 2 }
func (d itemDelegate) Spacing() int                              { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(PostItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s\n%s",
		lipgloss.NewStyle().Bold(true).Render(i.Title()),
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(i.Description()),
	)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type FeedModel struct {
	List   list.Model
	Loaded bool
}

func NewFeedModel() FeedModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Trending Posts"
	l.SetShowStatusBar(false)
	l.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	return FeedModel{
		List: l,
	}
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
