package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto/cli/internal/types"
)

type CommunityItem struct {
	types.Community
}

func (i CommunityItem) Title() string       { return "c/" + i.Community.Name }
func (i CommunityItem) Description() string { return fmt.Sprintf("%s • %d members", i.Community.Title, i.Community.MemberCount) }
func (i CommunityItem) FilterValue() string { return i.Community.Name + " " + i.Community.Title }

type CommunityModel struct {
	List   list.Model
	Loaded bool
}

func NewCommunityModel() CommunityModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Communities"
	l.SetShowStatusBar(false)
	l.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("34")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	return CommunityModel{
		List: l,
	}
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

func (m *CommunityModel) SetSize(width, height int) {
	m.List.SetSize(width, height)
}
