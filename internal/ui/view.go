package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"

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
	crumbs := []string{"DITTO"}
	switch m.State {
	case StateFeed:
		crumbs = append(crumbs, "FEED")
	case StatePostDetail:
		crumbs = append(crumbs, "c/"+m.PostDetailModel.Post.Community.Name, "p/"+m.PostDetailModel.Post.Title)
	case StateCommunities:
		crumbs = append(crumbs, "COMMUNITIES")
	case StateCreatePost:
		crumbs = append(crumbs, "NEW POST")
	case StateEditPost:
		crumbs = append(crumbs, "EDIT POST")
	case StateEditCommunity:
		crumbs = append(crumbs, "EDIT COMMUNITY")
	case StateProfileSettings:
		crumbs = append(crumbs, "SETTINGS")
	case StateHelp:
		crumbs = append(crumbs, "HELP")
	case StateLogin:
		crumbs = append(crumbs, "LOGIN")
	case StateRegister:
		crumbs = append(crumbs, "REGISTER")
	case StateSelection:
		crumbs = append(crumbs, "SELECT")
	case StateConfirm:
		crumbs = append(crumbs, "CONFIRM")
	}

	width := max(0, m.Width)
	left := m.renderBreadcrumbBar(crumbs)
	if width > 0 && lipgloss.Width(left) > width {
		left = m.renderBreadcrumbBar(m.compactBreadcrumbs(crumbs, width))
	}

	right := m.renderHeaderUserInfo()
	if right == "" {
		return left
	}

	if width == 0 {
		return lipgloss.JoinHorizontal(lipgloss.Left, left, " ", right)
	}

	if lipgloss.Width(left)+1+lipgloss.Width(right) <= width {
		return left + strings.Repeat(" ", width-lipgloss.Width(left)-lipgloss.Width(right)) + right
	}

	compactLeft := m.renderBreadcrumbBar(m.compactBreadcrumbs(crumbs, width))
	if lipgloss.Width(compactLeft)+1+lipgloss.Width(right) <= width {
		return compactLeft + strings.Repeat(" ", width-lipgloss.Width(compactLeft)-lipgloss.Width(right)) + right
	}

	if lipgloss.Width(right) > width {
		right = m.truncateHeaderLine(right, width)
	}
	return lipgloss.JoinVertical(lipgloss.Left, compactLeft, right)
}

func (m MainModel) renderBreadcrumbBar(crumbs []string) string {
	if len(crumbs) == 0 {
		return ""
	}

	primary := m.Theme.Title.Copy().Padding(0, 0)
	secondary := m.Theme.FooterSurface.Copy().Padding(0, 1).Foreground(m.Theme.TextSubtle.GetForeground())
	separator := m.Theme.TextSubtle.Render(" › ")

	parts := make([]string, 0, len(crumbs)*2-1)
	for i, crumb := range crumbs {
		label := "[" + m.truncateHeaderLabel(crumb, 28) + "]"
		style := secondary
		if i == 0 || i == len(crumbs)-1 {
			style = primary
		}
		parts = append(parts, style.Render(label))
		if i < len(crumbs)-1 {
			parts = append(parts, separator)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}

func (m MainModel) renderHeaderUserInfo() string {
	if m.Client.Token == "" {
		return ""
	}

	info := []string{"[logged in]"}
	if m.Me != nil {
		info[0] = fmt.Sprintf("[u/%s]", m.truncateHeaderLabel(m.Me.Username, 18))
	}
	if m.Wallet != nil {
		info = append(info, fmt.Sprintf("[🪙 %d]", m.Wallet.Coins), fmt.Sprintf("[💎 %d]", m.Wallet.Tokens))
	}

	pill := m.Theme.FooterSurface.Copy().Padding(0, 1).Foreground(m.Theme.Text.GetForeground())
	parts := make([]string, 0, len(info))
	for _, item := range info {
		parts = append(parts, pill.Render(item))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}

func (m MainModel) compactBreadcrumbs(crumbs []string, width int) []string {
	if len(crumbs) <= 2 {
		return crumbs
	}
	compact := []string{crumbs[0], "…", crumbs[len(crumbs)-1]}
	if width > 0 && lipgloss.Width(m.renderBreadcrumbBar(compact)) <= width {
		return compact
	}
	compact[2] = m.truncateHeaderLabel(compact[2], 20)
	return compact
}

func (m MainModel) truncateHeaderLine(s string, width int) string {
	if width <= 0 || lipgloss.Width(s) <= width {
		return s
	}
	plain := m.truncateHeaderLabel(stripHeaderBrackets(s), max(1, width))
	return m.Theme.TextSubtle.Render(plain)
}

func (m MainModel) truncateHeaderLabel(label string, limit int) string {
	if limit <= 0 || utf8.RuneCountInString(label) <= limit {
		return label
	}
	if limit == 1 {
		return "…"
	}
	return string([]rune(label)[:limit-1]) + "…"
}

func stripHeaderBrackets(s string) string {
	replacer := strings.NewReplacer("[", "", "]", "", "›", " ")
	return strings.Join(strings.Fields(replacer.Replace(s)), " ")
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
