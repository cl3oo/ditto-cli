package ui

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcku/ditto-cli/internal/api"
	"github.com/rfcku/ditto-cli/internal/config"
	"github.com/rfcku/ditto-cli/internal/types"
	"github.com/rfcku/ditto-cli/internal/ui/views"
)

func TestNewMainModel(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)

	if m.Client.BaseURL != cfg.BaseURL {
		t.Errorf("Expected baseURL %s, got %s", cfg.BaseURL, m.Client.BaseURL)
	}

	if m.State != StateLoading {
		t.Errorf("Expected initial state StateLoading, got %v", m.State)
	}
}

func TestMainModel_Update_ErrorUnauthorized(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)
	m.State = StateFeed

	// Simulate unauthorized error message
	newModel, _ := m.Update(errorMsg(api.ErrUnauthorized))
	updatedModel := newModel.(MainModel)

	if updatedModel.State != StateLogin {
		t.Errorf("Expected state StateLogin after ErrUnauthorized, got %v", updatedModel.State)
	}
}

func TestMainModel_Update_MeErrorUnauthorizedGoesToLogin(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)
	m.State = StateFeed
	m.Client.SetToken("still-valid-ish")

	newModel, _ := m.Update(meErrorMsg{err: api.ErrUnauthorized})
	updatedModel := newModel.(MainModel)

	if updatedModel.State != StateLogin {
		t.Fatalf("Expected state StateLogin after me unauthorized, got %v", updatedModel.State)
	}
	if updatedModel.Client.Token != "" {
		t.Fatalf("Expected token to be cleared after me unauthorized")
	}
}

func TestMainModel_Update_WalletErrorUnauthorizedIsIgnored(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)
	m.State = StateFeed
	m.Client.SetToken("still-valid-ish")
	m.Wallet = &types.Wallet{Coins: 42}

	newModel, _ := m.Update(walletErrorMsg{err: api.ErrUnauthorized})
	updatedModel := newModel.(MainModel)

	if updatedModel.State != StateFeed {
		t.Fatalf("Expected state to remain StateFeed after wallet unauthorized, got %v", updatedModel.State)
	}
	if updatedModel.Client.Token != "still-valid-ish" {
		t.Fatalf("Expected token to be preserved after wallet unauthorized")
	}
	if updatedModel.StatusMessage != "" {
		t.Fatalf("Expected no status message after wallet unauthorized, got %q", updatedModel.StatusMessage)
	}
	if updatedModel.Wallet != nil {
		t.Fatalf("Expected wallet to be cleared after wallet unauthorized")
	}
}

func TestMainModel_Update_LoginSuccess(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)
	m.State = StateLogin

	// Simulate login success message
	newModel, _ := m.Update(loginSuccessMsg("new-token"))
	updatedModel := newModel.(MainModel)

	if updatedModel.State != StateLogin {
		t.Errorf("Expected state StateLogin after loginSuccessMsg, got %v", updatedModel.State)
	}

	if updatedModel.Client.Token != "new-token" {
		t.Errorf("Expected client token to be updated")
	}
}

func TestMainModel_RenderFooterHelp(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)
	m.Width = 120

	tests := []struct {
		name     string
		state    State
		setup    func(*MainModel)
		contains []string
	}{
		{
			name:     "feed footer shows global and feed actions",
			state:    StateFeed,
			contains: []string{"Global", "Here", "n new post", "s search", ":man manual"},
		},
		{
			name:     "communities footer only shows truthful actions",
			state:    StateCommunities,
			contains: []string{"enter open", ":joined joined", ":random discover"},
		},
		{
			name:  "post detail reply mode shows reply hints",
			state: StatePostDetail,
			setup: func(m *MainModel) {
				m.PostDetailModel.SetShowCommentInput(true)
			},
			contains: []string{"type reply", "enter submit", "esc cancel"},
		},
		{
			name:  "post detail shows new reply shortcuts",
			state: StatePostDetail,
			setup: func(m *MainModel) {
				m.PostDetailModel.SelectedIdx = 0
			},
			contains: []string{"enter reply", "r reply", "p profile"},
		},
		{
			name:     "help footer shows manual navigation",
			state:    StateHelp,
			contains: []string{"j/k/gg/G scroll", "q close manual"},
		},
		{
			name:     "confirm footer shows confirm controls",
			state:    StateConfirm,
			contains: []string{"enter confirm", "q cancel", "esc cancel"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			local := m
			local.State = tt.state
			if tt.setup != nil {
				tt.setup(&local)
			}

			footer := local.renderFooterHelp()
			for _, want := range tt.contains {
				if !strings.Contains(footer, want) {
					t.Fatalf("footer %q missing %q", footer, want)
				}
			}
		})
	}
}

func TestMainModel_HelpManualContent(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)

	content := m.helpManualContent()
	for _, want := range []string{
		"Open the selected post or community",
		":joined**: List communities you joined.",
		fmt.Sprintf("**%s**: Open the command palette.", cfg.Keys.Palette),
		"reply to the post with **c**",
		"**Enter** to reply to a selected comment",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("help content %q missing %q", content, want)
		}
	}

	if strings.Contains(content, "join a community") {
		t.Fatalf("help content should not advertise join-on-enter")
	}
}

func TestMainModel_View_EmptyStates(t *testing.T) {
	cfg := config.DefaultConfig()

	tests := []struct {
		name     string
		state    State
		setup    func(*MainModel)
		contains []string
	}{
		{
			name:  "feed empty state shows actionable hints",
			state: StateFeed,
			setup: func(m *MainModel) {
				m.FeedModel.Loaded = true
				m.FeedModel.List.SetItems(nil)
			},
			contains: []string{"No posts yet", ":random discover posts", "n create a post"},
		},
		{
			name:  "communities empty state shows discovery hints",
			state: StateCommunities,
			setup: func(m *MainModel) {
				m.CommunityModel.Loaded = true
				m.CommunityModel.List.SetItems(nil)
			},
			contains: []string{"No communities to show", ":random discover communities", ":joined show joined communities"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMainModel(cfg)
			m.Width = 100
			m.Height = 30
			m.State = tt.state
			if tt.setup != nil {
				tt.setup(&m)
			}

			view := m.View()
			for _, want := range tt.contains {
				if !strings.Contains(view, want) {
					t.Fatalf("view %q missing %q", view, want)
				}
			}
		})
	}
}

func TestMainModel_RenderConfirmDialog(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)
	m.Width = 100
	m.Height = 30
	m.State = StateConfirm
	m.ConfirmDialog = ConfirmDialog{
		Action:      ConfirmDeleteAccount,
		Title:       "Delete account?",
		Body:        "This is irreversible.",
		TargetLabel: "tester",
		Step:        2,
		Steps:       2,
	}

	view := m.renderConfirmDialog()
	for _, want := range []string{"Confirm action (2/2)", "Delete account?", "Target: tester", "Enter to continue"} {
		if !strings.Contains(view, want) {
			t.Fatalf("confirm dialog %q missing %q", view, want)
		}
	}
}

func TestMainModel_Update_CommandPalette(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)
	m.State = StateFeed
	m.Width = 80
	m.Height = 24

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
	updatedModel := newModel.(MainModel)

	if updatedModel.State != StateCommandPalette {
		t.Errorf("Expected state StateCommandPalette, got %v", updatedModel.State)
	}

	if updatedModel.PreviousState != StateFeed {
		t.Errorf("Expected PreviousState StateFeed, got %v", updatedModel.PreviousState)
	}

	newModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updatedModel = newModel.(MainModel)

	if updatedModel.State != StateLoading {
		t.Errorf("Expected state StateLoading (fetching feed), got %v", updatedModel.State)
	}

	m.State = StateCommunities
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
	updatedModel = newModel.(MainModel)
	newModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updatedModel = newModel.(MainModel)

	if updatedModel.State != StateCommunities {
		t.Errorf("Expected state to return to StateCommunities, got %v", updatedModel.State)
	}
}

func TestMainModel_Update_PostDetailReplyShortcuts(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)
	m.State = StatePostDetail
	m.PostDetailModel.SetContent(
		types.Post{ID: "post-1", Title: "Test post", Author: types.UserMin{Username: "author"}},
		[]types.Comment{{ID: "comment-1", Content: "hello", Author: types.UserMin{Username: "replyguy"}}},
	)

	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	updatedModel := newModel.(MainModel)
	if !updatedModel.PostDetailModel.ShowCommentInput {
		t.Fatal("expected c to open reply input for the post")
	}
	if got := updatedModel.PostDetailModel.CommentInput.Placeholder; got != "Write a comment..." {
		t.Fatalf("expected top-level placeholder, got %q", got)
	}

	updatedModel.PostDetailModel.SetShowCommentInput(false)
	updatedModel.PostDetailModel.SelectedIdx = 0
	newModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updatedModel = newModel.(MainModel)
	if !updatedModel.PostDetailModel.ShowCommentInput {
		t.Fatal("expected enter to open reply input for the selected comment")
	}
	if got := updatedModel.PostDetailModel.CommentInput.Placeholder; got != "Replying to u/replyguy..." {
		t.Fatalf("expected comment reply placeholder, got %q", got)
	}
}

func TestMainModel_Update_EnterOnSearchPostOpensPostDetail(t *testing.T) {
	cfg := config.DefaultConfig()
	m := NewMainModel(cfg)
	m.State = StateCommunities
	m.CommunityModel.SetItems([]list.Item{views.PostItem{Post: types.Post{ID: "post-1", Title: "Search Result"}}})

	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updatedModel := newModel.(MainModel)

	if updatedModel.State != StateLoading {
		t.Fatalf("expected state %v, got %v", StateLoading, updatedModel.State)
	}
	if cmd == nil {
		t.Fatal("expected fetchPostDetail command")
	}
}

func TestMainModel_PerformSearch_ReturnsPartialResultsWhenPostSearchFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/search"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":[{"id":"c1","name":"golang","title":"Go","scores":{"sub_count":5,"post_count":2}}]}`))
		case strings.HasPrefix(r.URL.Path, "/posts"):
			http.Error(w, `{"error":"search unsupported"}`, http.StatusBadRequest)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	m := NewMainModel(config.DefaultConfig())
	m.Client = api.NewClient(server.URL)

	msg := m.performSearch("go")()
	search, ok := msg.(searchMsg)
	if !ok {
		t.Fatalf("expected searchMsg, got %T", msg)
	}
	if len(search) != 1 {
		t.Fatalf("expected 1 partial result, got %d", len(search))
	}
}

func TestMainModel_PerformSearch_ReturnsErrorWhenAllSearchesFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"broken"}`, http.StatusBadRequest)
	}))
	defer server.Close()

	m := NewMainModel(config.DefaultConfig())
	m.Client = api.NewClient(server.URL)

	msg := m.performSearch("go")()
	err, ok := msg.(errorMsg)
	if !ok {
		t.Fatalf("expected errorMsg, got %T", msg)
	}
	if !strings.Contains(err.Error(), "community search failed") || !strings.Contains(err.Error(), "post search failed") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
