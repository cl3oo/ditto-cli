package ui

import (
	"fmt"
	"strings"
)

func (m MainModel) globalFooterActions() []string {
	return []string{
		fmt.Sprintf("%s palette", m.Config.Keys.Palette),
		"enter open",
		"q back",
		":man manual",
	}
}

func (m MainModel) localFooterActions() []string {
	switch m.State {
	case StateFeed:
		return []string{
			"j/k move",
			fmt.Sprintf("%s/%s vote", m.Config.Keys.Upvote, m.Config.Keys.Downvote),
			fmt.Sprintf("%s new post", m.Config.Keys.New),
			fmt.Sprintf("%s refresh", m.Config.Keys.Refresh),
			"s search",
		}
	case StateCommunities:
		return []string{
			"j/k move",
			"enter open",
			fmt.Sprintf("%s new post", m.Config.Keys.New),
			":joined joined",
			":random discover",
		}
	case StatePostDetail:
		if m.PostDetailModel.ShowCommentInput {
			return []string{"type reply", "enter submit", "esc cancel"}
		}
		if m.PostDetailModel.SelectedIdx >= 0 {
			return []string{"j/k comments", "r reply", "p profile", "s share", "R report"}
		}
		return []string{"j/k comments", "r reply", "p author", "s share", "L load more"}
	case StateCreatePost:
		return []string{"tab next field", "shift+tab prev", "enter submit", "esc cancel"}
	case StateEditPost, StateEditCommunity, StateProfileSettings, StateRegister, StateLogin:
		return []string{"tab next field", "shift+tab prev", "enter submit", "esc cancel"}
	case StateHelp:
		return []string{"j/k scroll", "q close manual"}
	case StateSelection:
		return []string{"p post", "c community", "esc cancel"}
	case StateLoading:
		return []string{"wait", "q quit"}
	case StateConfirm:
		return []string{"enter confirm", "q cancel", "esc cancel"}
	case StateCommandPalette:
		return []string{"j/k move", "enter run", "esc cancel"}
	default:
		return nil
	}
}

func (m MainModel) helpManualContent() string {
	return strings.TrimSpace(fmt.Sprintf(`
# Ditto CLI Manual 📖

Welcome to Ditto. This guide sticks to what the TUI actually does right now.

## Starting the app 🚀
- Run **ditto-cli tui** to open the full terminal UI.
- Run **ditto-cli** without a subcommand when you want the root help and command list.

## Navigation 🕹️
- **j / k** or **Up / Down**: Move through lists and comments.
- **Enter**: Open the selected post or community, or submit the focused form.
- **q**: Go back to the previous screen.
- **esc**: Cancel command mode, dialogs, or inline reply/edit flows.
- **%s**: Open the command palette.

## Contextual actions by screen 🧭
- **Feed**: vote with **%s / %s**, create a post with **%s**, refresh with **%s**, search with **s**.
- **Communities**: open the selected community with **Enter**, discover joined communities with **:joined**, discover random communities with **:random**.
- **Post detail**: reply with **r**, inspect an author with **p**, share with **s**, report with **R**, and load more comments with **L**.
- **Forms** (login, register, compose, edit, settings): move with **Tab / Shift+Tab**, submit with **Enter**, cancel with **esc**.
- **Confirm dialogs**: confirm with **Enter**, cancel with **q** or **esc**.

## Global commands ⌨️
Press **:** to enter command mode:
- **:man**: Open this manual.
- **:feed**: Go to the feed.
- **:communities**: View communities.
- **:joined**: List communities you joined.
- **:new**: Create a new post.
- **:refresh**: Refresh the current view.
- **:random**: Load random posts or communities from the current list view.
- **:search <query>**: Search communities and posts.
- **:settings**: Edit your profile.
- **:logout**: Log out.
- **:q** or **:quit**: Exit the application.

## Moderation and advanced commands 🛡️
- **:delete**: Delete the current post or selected community when that screen supports it.
- **:edit**: Edit the current post or selected community when available.
- **:ban <userID>**: Ban a user from the current community.
- **:mod add <userID>**: Add a moderator.
- **:lock / :unlock**: Lock or unlock the current post.
- **:mod-delete <reason>**: Remove the current post as a moderator.

## Economy 💎
- **:award <awardID>**: Give an award to the current post.
- Your balance appears in the header when available.

---
Press **q** to return.
`, m.Config.Keys.Palette, m.Config.Keys.Upvote, m.Config.Keys.Downvote, m.Config.Keys.New, m.Config.Keys.Refresh))
}
