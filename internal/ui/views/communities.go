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

func (i CommunityItem) Title() string { return "c/" + i.Community.Name }
func (i CommunityItem) Description() string {
	return fmt.Sprintf("%s • 👥 %d • 📝 %d", i.Community.Title, i.Community.Scores.SubCount, i.Community.Scores.PostCount)
}
func (i CommunityItem) FilterValue() string { return i.Community.Name + " " + i.Community.Title }

type UserItem struct {
	types.User
}

func (i UserItem) Title() string       { return "u/" + i.User.Username }
func (i UserItem) Description() string { return "Followed user" }
func (i UserItem) FilterValue() string { return i.User.Username }

type communityDelegate struct {
	Theme lipgloss.Style
}

func (d communityDelegate) Height() int                               { return 2 }
func (d communityDelegate) Spacing() int                              { return 1 }
func (d communityDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d communityDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var titleStr, descStr string

	if i, ok := listItem.(CommunityItem); ok {
		titleStr = i.Title()
		descStr = i.Description()
	} else if i, ok := listItem.(UserItem); ok {
		titleStr = i.Title()
		descStr = i.Description()
	} else if i, ok := listItem.(PostItem); ok {
		titleStr = i.Title()
		descStr = i.Description()
	} else {
		return
	}

	title := lipgloss.NewStyle().Bold(true).Render(titleStr)
	desc := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(descStr)

	if index == m.Index() {
		title = d.Theme.Render(titleStr)
	}

	fmt.Fprintf(w, "  %s\n  %s", title, desc)
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
