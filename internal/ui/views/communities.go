package views

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/types"
)

type CommunityItem struct {
	types.Community
}

func (i CommunityItem) Title() string       { return "c/" + i.Community.Name }
func (i CommunityItem) Description() string { return i.Community.Title }
func (i CommunityItem) FilterValue() string { return i.Community.Name + " " + i.Community.Title }

type UserItem struct {
	types.User
}

func (i UserItem) Title() string       { return "u/" + i.User.Username }
func (i UserItem) Description() string { return "User profile" }
func (i UserItem) FilterValue() string { return i.User.Username }

type communityDelegate struct {
	Theme lipgloss.Style
}

func (d communityDelegate) Height() int                               { return 5 }
func (d communityDelegate) Spacing() int                              { return 1 }
func (d communityDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d communityDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	isSelected := index == m.Index()
	
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(0, 1).
		Width(m.Width() - 4)

	if isSelected {
		cardStyle = cardStyle.BorderForeground(lipgloss.Color("170"))
	}

	var header, title, stats string

	if i, ok := listItem.(CommunityItem); ok {
		header = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
			fmt.Sprintf("Community • Created %s", RelativeTime(i.CreatedAt)),
		)
		titleStyle := lipgloss.NewStyle().Bold(true)
		if isSelected {
			titleStyle = titleStyle.Foreground(lipgloss.Color("170"))
		}
		title = titleStyle.Render(i.Title() + ": " + i.Description())
		stats = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
			fmt.Sprintf("👥 %d members • 📝 %d posts", i.Scores.SubCount, i.Scores.PostCount),
		)
	} else if i, ok := listItem.(UserItem); ok {
		header = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
			fmt.Sprintf("User • Joined %s", RelativeTime(i.CreatedAt)),
		)
		titleStyle := lipgloss.NewStyle().Bold(true)
		if isSelected {
			titleStyle = titleStyle.Foreground(lipgloss.Color("170"))
		}
		title = titleStyle.Render(i.Title())
		stats = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("View profile for more details")
	} else if i, ok := listItem.(PostItem); ok {
		header = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
			fmt.Sprintf("u/%s in c/%s • %s", i.Author.Username, i.Community.Name, RelativeTime(i.CreatedAt)),
		)
		titleStyle := lipgloss.NewStyle().Bold(true)
		if isSelected {
			titleStyle = titleStyle.Foreground(lipgloss.Color("170"))
		}
		title = titleStyle.Render(i.Title())
		stats = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
			fmt.Sprintf("↑↓ %d • 💬 %d • 💎 %d", i.Scores.VoteScore, i.Scores.CommentCount, i.Scores.AwardCount),
		)
	} else {
		return
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		title,
		"",
		stats,
	)

	fmt.Fprint(w, cardStyle.Render(content))
}

type CommunityModel struct {
	List   list.Model
	Loaded bool
	Theme  lipgloss.Style
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

func (m *CommunityModel) SetTheme(selected lipgloss.Style) {
	m.Theme = selected
	m.List.SetDelegate(communityDelegate{Theme: selected})
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
