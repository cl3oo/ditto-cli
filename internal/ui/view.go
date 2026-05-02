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
		if m.FeedModel.Loaded && len(m.FeedModel.List.Items()) == 0 {
			content = m.renderEmptyState(
				"No posts yet",
				"The feed is empty right now. Try pulling in random posts, refreshing, or creating one.",
				[]string{":random discover posts", ":refresh reload feed", "n create a post", "s search"},
			)
		} else {
			content = m.FeedModel.View()
		}
	case StateLogin:
		content = m.LoginModel.View()
	case StatePostDetail:
		content = m.PostDetailModel.View()
	case StateCommunities:
		if m.CommunityModel.Loaded && len(m.CommunityModel.List.Items()) == 0 {
			content = m.renderEmptyState(
				"No communities to show",
				"There is nothing in this list yet. You can discover random communities or jump back to the feed.",
				[]string{":random discover communities", ":joined show joined communities", "q back to feed", "n create a post"},
			)
		} else {
			content = m.CommunityModel.View()
		}
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
	case StateConfirm:
		content = m.renderConfirmDialog()
	case StateCommandPalette:
		content = m.PaletteModel.View()
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
		style := m.Theme.SuccessBanner
		if strings.HasPrefix(m.StatusMessage, "Error") {
			style = m.Theme.ErrorBanner
		}
		s.WriteString("\n " + style.Render(m.StatusMessage))
	}

	if m.CommandBuffer != "" {
		s.WriteString("\n" + m.Theme.CommandSurface.Render(m.CommandBuffer))
	}

	if footer != "" {
		s.WriteString("\n" + footer)
	}

	return s.String()
}

func (m MainModel) renderFooterHelp() string {
	global := m.globalFooterActions()
	local := m.localFooterActions()

	section := func(title string, items []string, titleStyle, itemStyle lipgloss.Style) string {
		if len(items) == 0 {
			return ""
		}
		return lipgloss.JoinHorizontal(lipgloss.Left,
			titleStyle.Render(title),
			itemStyle.Render(strings.Join(items, "  •  ")),
		)
	}

	lineStyle := m.Theme.FooterSurface.Width(max(0, m.Width))

	return lipgloss.JoinVertical(lipgloss.Left,
		lineStyle.Render(section("Global", global, m.Theme.FooterSectionTitle, m.Theme.FooterItem)),
		lineStyle.Render(section("Here", local, m.Theme.FooterSectionTitleAlt, m.Theme.FooterItemStrong)),
	)
}

func (m MainModel) renderEmptyState(title, body string, actions []string) string {
	lines := []string{
		m.Theme.Title.Render(" " + title + " "),
		"",
		m.Theme.Text.Render(body),
	}
	if len(actions) > 0 {
		lines = append(lines, "", m.Theme.Text.Render("Try:"))
		for _, action := range actions {
			lines = append(lines, "  • "+m.Theme.TextSubtle.Render(action))
		}
	}

	box := m.Theme.Surface.
		BorderForeground(m.Theme.TextSubtle.GetForeground()).
		Width(min(max(52, m.Width-12), 88))

	return lipgloss.Place(
		max(0, m.Width),
		max(0, m.Height-6),
		lipgloss.Center,
		lipgloss.Center,
		box.Render(strings.Join(lines, "\n")),
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
	case StateConfirm:
		bc += "> CONFIRM "
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
		header += m.Theme.TextSubtle.
			MarginLeft(1).
			Render("(" + userStr + ")")
	}
	return header
}

func (m MainModel) renderConfirmDialog() string {
	stepLabel := ""
	if m.ConfirmDialog.Steps > 1 {
		stepLabel = fmt.Sprintf(" (%d/%d)", m.ConfirmDialog.Step, m.ConfirmDialog.Steps)
	}

	var lines []string
	lines = append(lines, m.Theme.Error.Bold(true).Render(" Confirm action"+stepLabel))
	lines = append(lines, "")
	lines = append(lines, m.Theme.Selected.Render(m.ConfirmDialog.Title))
	lines = append(lines, m.ConfirmDialog.Body)
	if m.ConfirmDialog.TargetLabel != "" {
		lines = append(lines, "")
		lines = append(lines, "Target: "+m.ConfirmDialog.TargetLabel)
	}
	if m.ConfirmDialog.Extra != "" {
		lines = append(lines, "Reason: "+m.ConfirmDialog.Extra)
	}
	lines = append(lines, "")
	lines = append(lines, "Enter to continue, q or esc to cancel.")

	box := m.Theme.ModalCritical.
		Width(min(max(60, m.Width-12), 90))

	return lipgloss.Place(
		max(0, m.Width),
		max(0, m.Height-6),
		lipgloss.Center,
		lipgloss.Center,
		box.Render(strings.Join(lines, "\n")),
	)
}
