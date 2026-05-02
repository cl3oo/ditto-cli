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

func TestGoldenFeedSelectionIndicator(t *testing.T) {
	m := newGoldenMainModel(82, 24, "")
	m.State = StateFeed
	m.FeedModel.Loaded = true
	m.FeedModel.SetPosts([]types.Post{
		{
			ID:        "post-1",
			Title:     "Arrow selection should be obvious without inverting the whole card",
			Author:    types.UserMin{Username: "guide"},
			Community: types.CommunityMin{Name: "ditto"},
			Scores:    types.Score{VoteScore: 12, CommentCount: 3, AwardCount: 1},
			CreatedAt: time.Now().Add(-90 * time.Minute),
		},
		{
			ID:        "post-2",
			Title:     "Unselected cards should stay calm",
			Author:    types.UserMin{Username: "rook"},
			Community: types.CommunityMin{Name: "ux"},
			Scores:    types.Score{VoteScore: 5, CommentCount: 1, AwardCount: 0},
			CreatedAt: time.Now().Add(-30 * time.Minute),
		},
	})
	m.FeedModel.SetSize(82, 16)

	assertGolden(t, "feed_selection_indicator.golden", m.FeedModel.View())
}

func TestGoldenCommunitySelectionIndicator(t *testing.T) {
	m := newGoldenMainModel(82, 24, "")
	m.State = StateCommunities
	m.CommunityModel.Loaded = true
	m.CommunityModel.SetCommunities([]types.Community{
		{
			ID:        "community-1",
			Name:      "ditto",
			Title:     "Core product",
			Scores:    types.Score{SubCount: 120, PostCount: 42},
			CreatedAt: time.Now().Add(-6 * time.Hour),
		},
		{
			ID:        "community-2",
			Name:      "design",
			Title:     "Design notes",
			Scores:    types.Score{SubCount: 24, PostCount: 9},
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
	})
	m.CommunityModel.SetSize(82, 16)

	assertGolden(t, "community_selection_indicator.golden", m.CommunityModel.View())
}

func TestGoldenLoginView(t *testing.T) {
	main := newGoldenMainModel(0, 0, "")
	m := views.NewLoginModel()
	m.SetTheme(main.Theme)
	m.Username.SetValue("astrocat")
	m.Password.SetValue("swordfish")
	m.Focused = 2
	m.Username.Blur()
	m.Password.Blur()

	assertGolden(t, "login_submit_focused.golden", m.View())
}

func TestGoldenRegisterViewWithError(t *testing.T) {
	main := newGoldenMainModel(0, 0, "")
	m := views.NewRegisterModel()
	m.SetTheme(main.Theme)
	m.Username.SetValue("astrocat")
	m.Email.SetValue("astrocat@example.com")
	m.Password.SetValue("swordfish")
	m.Confirm.SetValue("swordf1sh")
	m.Focused = 4
	m.Username.Blur()
	m.Email.Blur()
	m.Password.Blur()
	m.Confirm.Blur()
	m.Error = "passwords do not match"

	assertGolden(t, "register_error.golden", m.View())
}

func TestGoldenHelpViewNarrow(t *testing.T) {
	cfg := config.DefaultConfig()
	main := NewMainModel(cfg)
	m := views.NewHelpModel()
	m.SetContent(main.helpManualContent())
	m.SetSize(60, 18)

	assertGolden(t, "help_narrow.golden", m.View())
}

func TestGoldenPostDetailLargeContent(t *testing.T) {
	main := newGoldenMainModel(0, 0, "")
	m := views.NewPostDetailModel()
	m.SetTheme(main.Theme)
	m.SetSize(72, 22)

	now := time.Now()
	post := types.Post{
		ID:        "post-42",
		Title:     "Shipping snapshot coverage without making the TUI miserable",
		Content:   "This post body is intentionally long so the renderer has to wrap paragraphs across multiple lines.\n\nIt also includes a second paragraph to make sure markdown rendering stays readable in a constrained viewport.",
		Author:    types.UserMin{ID: "u-1", Username: "guide"},
		Community: types.CommunityMin{ID: "c-1", Name: "ditto"},
		Scores:    types.Score{VoteScore: 12, CommentCount: 3, AwardCount: 1},
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

func TestGoldenPostDetailNarrow(t *testing.T) {
	main := newGoldenMainModel(0, 0, "")
	m := views.NewPostDetailModel()
	m.SetTheme(main.Theme)
	m.SetSize(54, 18)

	now := time.Now()
	post := types.Post{
		ID:        "post-7",
		Title:     "Narrow layout should stay readable when metadata gets crowded",
		Content:   "A smaller viewport should still keep post detail usable without turning the header into soup.",
		Author:    types.UserMin{ID: "user-with-a-long-id", Username: "terminalfox"},
		Community: types.CommunityMin{ID: "community-with-a-long-id", Name: "ditto-design"},
		Scores:    types.Score{VoteScore: 19, CommentCount: 12, AwardCount: 2},
		CreatedAt: now.Add(-47 * time.Minute),
	}
	comments := []types.Comment{{
		ID:        "comment-1",
		Content:   "If this fits in a narrow fixture, it usually survives real terminals too.",
		Author:    types.UserMin{ID: "u-2", Username: "moss"},
		Scores:    types.Score{VoteScore: 6},
		CreatedAt: now.Add(-30 * time.Minute),
	}}

	m.SetContent(post, comments)
	m.SetSize(54, 18)

	assertGolden(t, "post_detail_narrow.golden", m.View())
}

func TestGoldenCommandPalette(t *testing.T) {
	m := newGoldenMainModel(80, 24, "")
	m.State = StateCommandPalette
	m.PaletteModel.SetTheme(m.Theme)
	m.PaletteModel.SetSize(80, 24)

	assertGolden(t, "command_palette.golden", m.View())
}

func TestGoldenConfirmDeleteAccount(t *testing.T) {
	m := newGoldenMainModel(84, 24, "")
	m.State = StateConfirm
	m.ConfirmDialog = ConfirmDialog{
		Action:      ConfirmDeleteAccount,
		Previous:    StateProfileSettings,
		Title:       "Delete account forever",
		Body:        "This permanently removes your account and cannot be undone.",
		TargetLabel: "u/astrocat",
		Extra:       "Deletes profile data, wallet access, and authored content associations.",
		Step:        2,
		Steps:       2,
	}

	assertGolden(t, "confirm_delete_account.golden", m.View())
}

func TestGoldenCreatePostNarrow(t *testing.T) {
	m := newGoldenMainModel(64, 22, "")
	m.State = StateCreatePost
	m.CreatePostModel.Title.SetValue("Patch narrow layout snapshots")
	m.CreatePostModel.CommunityID.SetValue("ditto")
	m.CreatePostModel.Content.SetValue("Golden fixtures should catch responsive regressions before users do.")
	m.CreatePostModel.Focused = 2
	m.CreatePostModel.Title.Blur()
	m.CreatePostModel.CommunityID.Blur()
	m.CreatePostModel.Content.Focus()
	m.CreatePostModel.SetTheme(m.Theme)

	assertGolden(t, "create_post_narrow.golden", m.View())
}

func TestGoldenSettingsLightTheme(t *testing.T) {
	m := newGoldenMainModel(76, 20, "light")
	m.State = StateProfileSettings
	m.Client.SetToken("token-123")
	m.Me = &types.User{Username: "sunny"}
	m.Wallet = &types.Wallet{Coins: 42, Tokens: 7}
	m.SettingsModel.Avatar.SetValue("https://cdn.example.com/avatar.png")
	m.SettingsModel.Focused = 1
	m.SettingsModel.Avatar.Blur()
	m.SettingsModel.SetTheme(m.Theme)

	assertGolden(t, "settings_light_theme.golden", m.View())
}

func newGoldenMainModel(width, height int, themeName string) MainModel {
	cfg := config.DefaultConfig()
	if themeName != "" {
		cfg.Appearance.Theme = themeName
	}
	m := NewMainModel(cfg)
	m.Width = width
	m.Height = height
	return m
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
