package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rfcku/ditto-cli/internal/types"
	"github.com/rfcku/ditto-cli/internal/ui/theme"
)

type listCard struct {
	KindLabel string
	Meta      []string
	Title     string
	Body      string
	Stats     []string
}

func newPostListCard(post types.Post, width int) listCard {
	return listCard{
		KindLabel: "POST",
		Meta: []string{
			fmt.Sprintf("c/%s", post.Community.Name),
			fmt.Sprintf("u/%s", post.Author.Username),
			RelativeTime(post.CreatedAt),
		},
		Title: post.Title,
		Body:  clampLine(post.Content, max(24, width-14)),
		Stats: []string{
			fmt.Sprintf("↑↓ %d", post.Scores.VoteScore),
			fmt.Sprintf("💬 %d", post.Scores.CommentCount),
			fmt.Sprintf("💎 %d", post.Scores.AwardCount),
		},
	}
}

func newCommunityListCard(community types.Community, width int) listCard {
	return listCard{
		KindLabel: "COMMUNITY",
		Meta: []string{
			fmt.Sprintf("c/%s", community.Name),
			RelativeTime(community.CreatedAt),
		},
		Title: community.Title,
		Body:  clampLine(community.Description, max(24, width-14)),
		Stats: []string{
			fmt.Sprintf("👥 %d members", community.Scores.SubCount),
			fmt.Sprintf("📝 %d posts", community.Scores.PostCount),
		},
	}
}

func newUserListCard(user types.User) listCard {
	return listCard{
		KindLabel: "USER",
		Meta: []string{
			fmt.Sprintf("u/%s", user.Username),
			fmt.Sprintf("joined %s", RelativeTime(user.CreatedAt)),
		},
		Title: fmt.Sprintf("u/%s", user.Username),
		Body:  "Ditto user profile",
		Stats: []string{"Followed user"},
	}
}

func renderListCard(cardStyle lipgloss.Style, isSelected bool, t theme.Theme, card listCard) string {
	titleStyle := t.Text.Bold(true)
	if isSelected {
		titleStyle = t.AccentText.Bold(true)
	}

	metaParts := append([]string{card.KindLabel}, card.Meta...)
	meta := strings.Join(filterEmpty(metaParts), " • ")
	stats := strings.Join(filterEmpty(card.Stats), " • ")

	lines := []string{t.TextSubtle.Render(meta), titleStyle.Render(card.Title)}
	if card.Body != "" {
		lines = append(lines, t.Text.Render(card.Body))
	}
	if stats != "" {
		lines = append(lines, t.TextSubtle.Render(stats))
	}

	return withSelectionIndicator(cardStyle.Render(strings.Join(lines, "\n\n")), isSelected, t)
}

func renderVerticalCard(cardStyle lipgloss.Style, isSelected bool, t theme.Theme, meta, title, body, stats string) string {
	return renderListCard(cardStyle, isSelected, t, listCard{
		Meta:  []string{meta},
		Title: title,
		Body:  body,
		Stats: []string{stats},
	})
}

func clampLine(text string, limit int) string {
	text = strings.TrimSpace(strings.Join(strings.Fields(text), " "))
	if text == "" || limit <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	if limit == 1 {
		return "…"
	}
	return strings.TrimSpace(string(runes[:limit-1])) + "…"
}

func filterEmpty(parts []string) []string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}
		filtered = append(filtered, part)
	}
	return filtered
}
