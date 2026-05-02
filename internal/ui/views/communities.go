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

	var card listCard

	if i, ok := listItem.(CommunityItem); ok {
		card = newCommunityListCard(i.Community, m.Width())
	} else if i, ok := listItem.(UserItem); ok {
		card = newUserListCard(i.User)
	} else if i, ok := listItem.(PostItem); ok {
		card = newPostListCard(i.Post, m.Width())
	} else {
		return
	}

	_, _ = fmt.Fprint(w, renderListCard(cardStyle, isSelected, d.Theme, card))
}

type CommunityModel struct {
	List     list.Model
	Loaded   bool
	Theme    theme.Theme
	pendingG bool
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "g":
			if m.pendingG {
				m.List.Select(0)
				m.pendingG = false
				return m, nil
			}
			m.pendingG = true
			return m, nil
		case "G":
			m.List.Select(len(m.List.Items()) - 1)
			m.pendingG = false
			return m, nil
		case "H":
			m.List.Select(m.List.Paginator.Page * m.List.Paginator.PerPage)
			m.pendingG = false
			return m, nil
		case "L":
			last := (m.List.Paginator.Page+1)*m.List.Paginator.PerPage - 1
			if last >= len(m.List.Items()) {
				last = len(m.List.Items()) - 1
			}
			m.List.Select(last)
			m.pendingG = false
			return m, nil
		case "M":
			start := m.List.Paginator.Page * m.List.Paginator.PerPage
			last := (m.List.Paginator.Page+1)*m.List.Paginator.PerPage - 1
			if last >= len(m.List.Items()) {
				last = len(m.List.Items()) - 1
			}
			m.List.Select(start + (last-start)/2)
			m.pendingG = false
			return m, nil
		default:
			m.pendingG = false
		}
	}

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
