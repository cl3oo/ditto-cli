package views

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcku/ditto-cli/internal/types"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type CommunityItem struct {
	types.Community
}

func (i CommunityItem) Title() string { return "c/" + i.Name }
func (i CommunityItem) Description() string {
	return fmt.Sprintf("%s • 👥 %d • 📝 %d", i.Community.Title, i.Scores.SubCount, i.Scores.PostCount)
}
func (i CommunityItem) FilterValue() string { return i.Name + " " + i.Community.Title }

type UserItem struct {
	types.User
}

func (i UserItem) Title() string       { return "u/" + i.Username }
func (i UserItem) Description() string { return "Followed user" }
func (i UserItem) FilterValue() string { return i.Username }

type communityDelegate struct {
	Theme theme.Theme
}

func (d communityDelegate) Height() int                               { return 5 }
func (d communityDelegate) Spacing() int                              { return 1 }
func (d communityDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d communityDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	isSelected := index == m.Index()

	cardStyle := d.Theme.Post.
		Width(max(20, m.Width()-6))

	var meta, title, body, stats string

	if i, ok := listItem.(CommunityItem); ok {
		meta = fmt.Sprintf("c/%s • created %s", i.Name, RelativeTime(i.CreatedAt))
		title = i.Community.Title
		body = clampLine(i.Community.Description, max(24, m.Width()-14))
		stats = fmt.Sprintf("👥 %d members • 📝 %d posts", i.Scores.SubCount, i.Scores.PostCount)
	} else if i, ok := listItem.(UserItem); ok {
		meta = fmt.Sprintf("u/%s • joined %s", i.Username, RelativeTime(i.CreatedAt))
		title = i.Title()
		body = "Ditto user"
		stats = i.Description()
	} else if i, ok := listItem.(PostItem); ok {
		meta = fmt.Sprintf("c/%s • u/%s • %s", i.Community.Name, i.Author.Username, RelativeTime(i.CreatedAt))
		title = i.Title()
		body = clampLine(i.Content, max(24, m.Width()-14))
		stats = fmt.Sprintf("↑↓ %d • 💬 %d • 💎 %d", i.Scores.VoteScore, i.Scores.CommentCount, i.Scores.AwardCount)
	} else {
		return
	}

	_, _ = fmt.Fprint(w, renderVerticalCard(cardStyle, isSelected, d.Theme, meta, title, body, stats))
}

type CommunityModel struct {
	List   list.Model
	Loaded bool
	Theme  theme.Theme
}

func NewCommunityModel() CommunityModel {
	l := list.New([]list.Item{}, communityDelegate{}, 0, 0)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return CommunityModel{
		List: l,
	}
}

func (m *CommunityModel) SetTheme(t theme.Theme) {
	m.Theme = t
	m.List.SetDelegate(communityDelegate{Theme: t})
}
func (m CommunityModel) Init() tea.Cmd {
	return nil
}

func (m CommunityModel) Update(msg tea.Msg) (CommunityModel, tea.Cmd) {
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m CommunityModel) View() string {
	return m.List.View()
}

func (m *CommunityModel) SetCommunities(communities []types.Community) {
	items := make([]list.Item, len(communities))
	for i, c := range communities {
		items[i] = CommunityItem{Community: c}
	}
	m.List.SetItems(items)
	m.Loaded = true
}

func (m *CommunityModel) SetUsers(users []types.User) {
	items := make([]list.Item, len(users))
	for i, u := range users {
		items[i] = UserItem{User: u}
	}
	m.List.SetItems(items)
	m.Loaded = true
}

func (m *CommunityModel) SetItems(items []list.Item) {
	m.List.SetItems(items)
	m.Loaded = true
}

func (m *CommunityModel) SetSize(width, height int) {
	m.List.SetSize(width, height)
}
