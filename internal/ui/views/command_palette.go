package views

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
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
	Theme theme.Theme
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
		str = d.Theme.Selected.Render("> " + i.TitleText)
		desc = d.Theme.TextSubtle.Render(i.DescriptionText)
	} else {
		str = lipgloss.NewStyle().PaddingLeft(2).Render(i.TitleText)
		desc = d.Theme.TextSubtle.PaddingLeft(2).Render(i.DescriptionText)
	}

	_, _ = fmt.Fprintf(w, "%s\n%s", str, desc)
}

type CommandPaletteModel struct {
	List  list.Model
	Theme theme.Theme
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

	l := list.New(items, commandDelegate{}, 0, 0)
	l.Title = "Command Palette"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return CommandPaletteModel{
		List: l,
	}
}

func (m *CommandPaletteModel) SetTheme(t theme.Theme) {
	m.Theme = t
	m.List.SetDelegate(commandDelegate{Theme: t})
	m.List.Styles.Title = t.Title
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
