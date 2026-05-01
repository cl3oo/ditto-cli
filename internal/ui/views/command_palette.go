package views

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CommandAction string

const (
	ActionGoFeed        CommandAction = "go-feed"
	ActionGoCommunities CommandAction = "go-communities"
	ActionGoLogin       CommandAction = "go-login"
	ActionGoRegister    CommandAction = "go-register"
	ActionCreatePost    CommandAction = "create-post"
	ActionGoSettings    CommandAction = "go-settings"
	ActionGoHelp        CommandAction = "go-help"
	ActionQuit          CommandAction = "quit"
)

type CommandItem struct {
	Action          CommandAction
	TitleText       string
	DescriptionText string
}

func (i CommandItem) Title() string       { return i.TitleText }
func (i CommandItem) Description() string { return i.DescriptionText }
func (i CommandItem) FilterValue() string { return i.TitleText + " " + i.DescriptionText }

type commandDelegate struct {
	ActiveStyle lipgloss.Style
}

func (d commandDelegate) Height() int                               { return 2 }
func (d commandDelegate) Spacing() int                              { return 0 }
func (d commandDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d commandDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(CommandItem)
	if !ok {
		return
	}

	var str, desc string
	if index == m.Index() {
		str = d.ActiveStyle.Render("> " + i.TitleText)
		desc = d.ActiveStyle.Foreground(lipgloss.Color("241")).Render(i.DescriptionText)
	} else {
		str = lipgloss.NewStyle().PaddingLeft(2).Render(i.TitleText)
		desc = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("241")).Render(i.DescriptionText)
	}

	_, _ = fmt.Fprintf(w, "%s\n%s", str, desc)
}

type CommandPaletteModel struct {
	List  list.Model
	Theme lipgloss.Style
}

func NewCommandPaletteModel() CommandPaletteModel {
	items := []list.Item{
		CommandItem{Action: ActionGoFeed, TitleText: "Feed", DescriptionText: "View home feed"},
		CommandItem{Action: ActionGoCommunities, TitleText: "Communities", DescriptionText: "Browse all communities"},
		CommandItem{Action: ActionGoLogin, TitleText: "Login", DescriptionText: "Sign in to your account"},
		CommandItem{Action: ActionGoRegister, TitleText: "Register", DescriptionText: "Create a new account"},
		CommandItem{Action: ActionCreatePost, TitleText: "Create Post", DescriptionText: "Share something new"},
		CommandItem{Action: ActionGoSettings, TitleText: "Settings", DescriptionText: "Configure app preferences"},
		CommandItem{Action: ActionGoHelp, TitleText: "Help", DescriptionText: "Show keyboard shortcuts"},
		CommandItem{Action: ActionQuit, TitleText: "Quit", DescriptionText: "Exit the application"},
	}

	l := list.New(items, commandDelegate{ActiveStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("170"))}, 0, 0)
	l.Title = "Command Palette"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("170")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	return CommandPaletteModel{
		List: l,
	}
}

func (m *CommandPaletteModel) SetTheme(selected lipgloss.Style) {
	m.Theme = selected
	m.List.SetDelegate(commandDelegate{ActiveStyle: selected})
	m.List.Styles.Title = selected.
		Background(selected.GetForeground()).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)
}

func (m CommandPaletteModel) Update(msg tea.Msg) (CommandPaletteModel, tea.Cmd) {
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m CommandPaletteModel) View() string {
	return m.List.View()
}

func (m *CommandPaletteModel) SetSize(width, height int) {
	m.List.SetSize(width, height)
}
