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
	case StateRegister:
		content = m.RegisterModel.View()
	case StateSelection:
		content = m.renderSelectionMenu()
	}
	s.WriteString(content)

	// Command Buffer and Status (pinned to bottom)
	footerHeight := 0
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

	// Status Message
	if m.StatusMessage != "" {
		style := m.Theme.Selected
		if strings.HasPrefix(m.StatusMessage, "Error") {
			style = m.Theme.Error
		}
		s.WriteString("\n " + style.Render(m.StatusMessage))
	}

	// Command Buffer
	if m.CommandBuffer != "" {
		s.WriteString("\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(m.Theme.Accent).
			Padding(0, 1).
			Render(m.CommandBuffer))
	}

	return s.String()
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
			userStr = fmt.Sprintf("(u/%s)", m.Me.Username)
		}
		header += lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(1).
			Render(userStr)
	}
	return header
}
