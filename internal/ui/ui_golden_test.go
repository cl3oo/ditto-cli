package ui

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/rfcku/ditto-cli/internal/config"
	"github.com/rfcku/ditto-cli/internal/types"
	"github.com/rfcku/ditto-cli/internal/ui/views"
)

var updateGolden = flag.Bool("update", false, "update golden files")

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestGoldenMainViewFeedEmpty(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)
	m.Width = 90
	m.Height = 26
	m.State = StateFeed
	m.FeedModel.Loaded = true
	m.FeedModel.List.SetItems(nil)

	assertGolden(t, "main_feed_empty.golden", m.View())
}

func TestGoldenMainViewCommunitiesEmpty(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)
	m.Width = 90
	m.Height = 26
	m.State = StateCommunities
	m.CommunityModel.Loaded = true
	m.CommunityModel.List.SetItems(nil)

	assertGolden(t, "main_communities_empty.golden", m.View())
}

func TestGoldenLoginView(t *testing.T) {
	m := views.NewLoginModel()
	m.Username.SetValue("astrocat")
	m.Password.SetValue("swordfish")
	m.Focused = 2
	m.Username.Blur()
	m.Password.Blur()

	assertGolden(t, "login_submit_focused.golden", m.View())
}

func TestGoldenHelpViewNarrow(t *testing.T) {
	m := views.NewHelpModel()
	m.SetSize(60, 18)
	m.Ready = true
	m.SetSize(60, 18)

	assertGolden(t, "help_narrow.golden", m.View())
}

func TestGoldenPostDetailLargeContent(t *testing.T) {
	cfg := config.DefaultConfig()
	main := NewMainModel(cfg)
	m := views.NewPostDetailModel()
	m.SetTheme(main.Theme.Accent, main.Theme.Selected, main.Theme.Markdown)
	m.SetSize(72, 22)

	now := time.Now()
	post := types.Post{
		ID:      "post-42",
		Title:   "Shipping snapshot coverage without making the TUI miserable",
		Content: "This post body is intentionally long so the renderer has to wrap paragraphs across multiple lines.\n\nIt also includes a second paragraph to make sure markdown rendering stays readable in a constrained viewport.",
		Author:  types.UserMin{ID: "u-1", Username: "guide"},
		Community: types.CommunityMin{ID: "c-1", Name: "ditto"},
		Scores: types.Score{VoteScore: 12, CommentCount: 3, AwardCount: 1},
		CreatedAt: now.Add(-2 * time.Hour),
	}
	comments := []types.Comment{
		{
			ID:        "c-1",
			Content:   "Love this direction. Golden tests are boring right until they save a release.",
			Author:    types.UserMin{ID: "u-2", Username: "sable"},
			Scores:    types.Score{VoteScore: 7},
			CreatedAt: now.Add(-95 * time.Minute),
			Children: []types.Comment{{
				ID:        "c-2",
				Content:   "Yep. They are basically guard rails with better manners.",
				Author:    types.UserMin{ID: "u-3", Username: "rook"},
				Scores:    types.Score{VoteScore: 4},
				CreatedAt: now.Add(-70 * time.Minute),
			}},
		},
		{
			ID:        "c-3",
			Content:   "Please include a narrow-width fixture too, terminals love chaos.",
			Author:    types.UserMin{ID: "u-4", Username: "moss"},
			Scores:    types.Score{VoteScore: 5},
			CreatedAt: now.Add(-35 * time.Minute),
		},
	}

	m.SetContent(post, comments)
	m.SelectedIdx = 1
	m.SetSize(72, 22)

	assertGolden(t, "post_detail_large_content.golden", m.View())
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", name)
	normalized := normalizeGolden(got)

	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir golden dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(normalized), 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
	}

	wantBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v", name, err)
	}

	want := string(wantBytes)
	if normalized != want {
		t.Fatalf("golden mismatch for %s\n\n--- want ---\n%s\n--- got ---\n%s", name, want, normalized)
	}
}

func normalizeGolden(s string) string {
	s = ansiPattern.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n") + "\n"
}
