package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m MainModel) View() string {
	var s strings.Builder

	// Header
	header := HeaderStyle.Render(" DITTO CLI ")
	if m.Client.Token != "" {
		header += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(" (Logged In)")
	}
	s.WriteString(header + "\n\n")

	if m.Error != nil {
		s.WriteString(ErrStyle.Render(fmt.Sprintf("Error: %v", m.Error)))
		s.WriteString("\n\n(r to retry, L for login)")
	} else {
		switch m.State {
		case StateLoading:
			s.WriteString(" Loading...")
		case StateFeed:
			s.WriteString(m.FeedModel.View())
		case StateLogin:
			s.WriteString(m.LoginModel.View())
		case StatePostDetail:
			s.WriteString(m.PostDetailModel.View())
		case StateCommunities:
			s.WriteString(m.CommunityModel.View())
		case StateCreatePost:
			s.WriteString(m.CreatePostModel.View())
		}
	}

	// Footer
	var footer string
	if m.State == StatePostDetail {
		footer = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("q: quit • esc/backspace: back • ↑/↓: scroll")
	} else if m.State == StateCommunities {
		footer = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("q: quit • esc/backspace: back • F: feed • r: refresh")
	} else if m.State == StateCreatePost {
		footer = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("q: quit • esc/backspace: back • tab: next field")
	} else {
		footer = "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("q: quit • L: login • F: feed • C: communities • n: new post • r: refresh • enter: view post")
	}
	s.WriteString(footer)

	return s.String()
}
