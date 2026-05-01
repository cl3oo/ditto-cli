package views

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type HelpModel struct {
	Viewport viewport.Model
	Ready    bool
	Width    int
	Height   int
	Theme    theme.Theme
}

func NewHelpModel() HelpModel {
	return HelpModel{
		Viewport: viewport.New(0, 0),
	}
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	var cmd tea.Cmd
	m.Viewport, cmd = m.Viewport.Update(msg)
	return m, cmd
}

func (m HelpModel) View() string {
	if !m.Ready {
		return "  Loading manual..."
	}
	return m.Viewport.View()
}

func (m *HelpModel) SetSize(width, height int) {
	m.Width = width
	m.Height = height
	m.Viewport.Width = width
	m.Viewport.Height = height
	m.render()
}

func (m *HelpModel) SetTheme(t theme.Theme) {
	m.Theme = t
	if m.Ready {
		m.render()
	}
}

func (m *HelpModel) render() {
	if m.Width == 0 || m.Height == 0 {
		return
	}

	content := `
# Ditto CLI Manual 📖

Welcome to Ditto! This guide will help you navigate and interact with the platform.

## Navigation 🕹️
- **j / k** or **Up / Down**: Navigate through lists (Feed, Communities, Comments).
- **Enter**: Open a post, join a community, or submit a form.
- **q**: Go back to the previous screen.
- **esc**: Clear the command buffer or exit a sub-menu.

## Post Detail Actions 📄
- **j / k**: Select specific comments.
- **r**: Reply to the selected comment (or the post if nothing is selected).
- **p**: View the profile of the comment author.
- **:cc**: Write a top-level comment.

## Global Commands ⌨️
Press **:** to enter command mode:
- **:feed**: Go to the trending feed.
- **:communities**: View all communities.
- **:new**: Create a new post.
- **:search <query>**: Search for communities and posts.
- **:settings**: Edit your profile (avatar).
- **:logout**: Log out of your account.
- **:q** or **:quit**: Exit the application.

## Moderation (Moderators Only) 🛡️
- **:ban <userID>**: Ban a user from the community.
- **:mod add <userID>**: Add a new moderator.
- **:lock / :unlock**: Lock or unlock the current post.
- **:mod-delete <reason>**: Remove a post.

## Economy 💎
- **:award <awardID>**: Give an award to the current post.
- Your balance is displayed in the header.

---
Press **q** to return to the feed.
`
	var renderer *glamour.TermRenderer
	if m.Theme.Markdown != "" {
		renderer, _ = glamour.NewTermRenderer(
			glamour.WithStylesFromJSONBytes([]byte(m.Theme.Markdown)),
			glamour.WithWordWrap(m.Width-4),
		)
	} else {
		renderer, _ = glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(m.Width-4),
		)
	}
	out, _ := renderer.Render(content)
	m.Viewport.SetContent(out)
	m.Ready = true
}
