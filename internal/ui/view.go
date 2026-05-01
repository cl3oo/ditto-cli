package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m MainModel) View() string {
	var s strings.Builder

	// Header
	s.WriteString(m.renderHeader() + "\n\n")

	// Main Content
	var content string
	switch m.State {
	case StateLoading:
		content = "  Loading..."
	case StateFeed:
		content = m.FeedModel.View()
	case StateLogin:
		content = m.LoginModel.View()
	case StatePostDetail:
		content = m.PostDetailModel.View()
	case StateCommunities:
		content = m.CommunityModel.View()
	case StateCreatePost:
		content = m.CreatePostModel.View()
	case StateEditPost:
		content = m.EditPostModel.View()
	case StateEditCommunity:
		content = m.EditCommunityModel.View()
	case StateProfileSettings:
		content = m.SettingsModel.View()
	case StateHelp:
		content = m.HelpModel.View()
	case StateRegister:
		content = m.RegisterModel.View()
	case StateSelection:
		content = m.renderSelectionMenu()
	}
	s.WriteString(content)

	footer := m.renderFooterHelp()
	footerHeight := lipgloss.Height(footer)
	if m.CommandBuffer != "" {
		footerHeight++
	}
	if m.StatusMessage != "" {
		footerHeight++
	}

	contentHeight := lipgloss.Height(content)
	headerHeight := lipgloss.Height(m.renderHeader()) + 2

	padding := m.Height - contentHeight - headerHeight - footerHeight
	if padding > 0 {
		s.WriteString(strings.Repeat("\n", padding))
	}

	if m.StatusMessage != "" {
		style := m.Theme.Selected
		if strings.HasPrefix(m.StatusMessage, "Error") {
			style = m.Theme.Error
		}
		s.WriteString("\n " + style.Render(m.StatusMessage))
	}

	if m.CommandBuffer != "" {
		s.WriteString("\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(m.Theme.Accent).
			Padding(0, 1).
			Render(m.CommandBuffer))
	}

	if footer != "" {
		s.WriteString("\n" + footer)
	}

	return s.String()
}

func (m MainModel) renderFooterHelp() string {
	global := []string{"/: cmd", "enter open", "q back", "? :man"}
	local := []string{}

	switch m.State {
	case StateFeed:
		local = []string{
			"j/k move",
			fmt.Sprintf("%s/%s vote", m.Config.Keys.Upvote, m.Config.Keys.Downvote),
			fmt.Sprintf("%s new post", m.Config.Keys.New),
			fmt.Sprintf("%s refresh", m.Config.Keys.Refresh),
			"s search",
		}
	case StateCommunities:
		local = []string{"j/k move", "enter open/join", fmt.Sprintf("%s new post", m.Config.Keys.New), ":random discover"}
	case StatePostDetail:
		if m.PostDetailModel.ShowCommentInput {
			local = []string{"type reply", "enter submit", "esc cancel"}
		} else if m.PostDetailModel.SelectedIdx >= 0 {
			local = []string{"j/k comments", "r reply", "p profile", "s share", "R report"}
		} else {
			local = []string{"j/k comments", "r reply", "p author", "s share", "L load more"}
		}
	case StateCreatePost:
		local = []string{"tab next field", "shift+tab prev", "enter submit", "esc cancel"}
	case StateEditPost, StateEditCommunity, StateProfileSettings, StateRegister, StateLogin:
		local = []string{"tab next field", "shift+tab prev", "enter submit", "esc cancel"}
	case StateHelp:
		local = []string{"j/k scroll", "q close manual"}
	case StateSelection:
		local = []string{"p post", "c community", "esc cancel"}
	case StateLoading:
		local = []string{"wait", "q quit"}
	}

	section := func(title string, items []string, titleStyle, itemStyle lipgloss.Style) string {
		if len(items) == 0 {
			return ""
		}
		return lipgloss.JoinHorizontal(lipgloss.Left,
			titleStyle.Render(title),
			itemStyle.Render(strings.Join(items, "  •  ")),
		)
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Padding(0, 1)
	globalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Padding(0, 1)
	localTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(m.Theme.Accent).Padding(0, 1)
	localStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Padding(0, 1)
	lineStyle := lipgloss.NewStyle().Background(lipgloss.Color("236")).Width(max(0, m.Width))

	return lipgloss.JoinVertical(lipgloss.Left,
		lineStyle.Render(section("Global", global, titleStyle, globalStyle)),
		lineStyle.Render(section("Here", local, localTitleStyle, localStyle)),
	)
}

func (m MainModel) renderSelectionMenu() string {
	var s strings.Builder
	title := m.Theme.Title.Render(" Create New ")
	s.WriteString(" " + title + "\n\n")
	s.WriteString("  [p] Post\n")
	s.WriteString("  [c] Community\n")
	s.WriteString("\n  Press p or c to select, esc to cancel")
	return s.String()
}

func (m MainModel) renderHeader() string {
	bc := " DITTO "
	switch m.State {
	case StateFeed:
		bc += "> FEED "
	case StatePostDetail:
		bc += fmt.Sprintf("> c/%s > p/%s ", m.PostDetailModel.Post.Community.Name, m.PostDetailModel.Post.Title)
	case StateCommunities:
		bc += "> COMMUNITIES "
	case StateCreatePost:
		bc += "> NEW POST "
	case StateEditPost:
		bc += "> EDIT POST "
	case StateEditCommunity:
		bc += "> EDIT COMMUNITY "
	case StateProfileSettings:
		bc += "> SETTINGS "
	case StateHelp:
		bc += "> HELP "
	case StateLogin:
		bc += "> LOGIN "
	case StateRegister:
		bc += "> REGISTER "
	case StateSelection:
		bc += "> SELECT "
	}

	header := m.Theme.Title.Render(bc)
	if m.Client.Token != "" {
		userStr := "(Logged In)"
		if m.Me != nil {
			userStr = fmt.Sprintf("u/%s", m.Me.Username)
		}
		if m.Wallet != nil {
			userStr += fmt.Sprintf(" | 🪙 %d | 💎 %d", m.Wallet.Coins, m.Wallet.Tokens)
		}
		header += lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(1).
			Render("(" + userStr + ")")
	}
	return header
}
