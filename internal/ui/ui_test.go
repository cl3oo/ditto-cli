package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcku/ditto-cli/internal/api"
	"github.com/rfcku/ditto-cli/internal/config"
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
			contains: []string{"Global", "Here", "n new post", "s search"},
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
			name:     "help footer shows manual navigation",
			state:    StateHelp,
			contains: []string{"j/k scroll", "q close manual"},
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

	// 1. Open palette with ":"
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
	updatedModel := newModel.(MainModel)

	if updatedModel.State != StateCommandPalette {
		t.Errorf("Expected state StateCommandPalette, got %v", updatedModel.State)
	}

	if updatedModel.PreviousState != StateFeed {
		t.Errorf("Expected PreviousState StateFeed, got %v", updatedModel.PreviousState)
	}

	// 2. Select Feed (should already be first item) and press enter
	newModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updatedModel = newModel.(MainModel)

	if updatedModel.State != StateLoading {
		t.Errorf("Expected state StateLoading (fetching feed), got %v", updatedModel.State)
	}

	// 3. Test escape
	m.State = StateCommunities
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(":")})
	updatedModel = newModel.(MainModel)
	newModel, _ = updatedModel.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updatedModel = newModel.(MainModel)

	if updatedModel.State != StateCommunities {
		t.Errorf("Expected state to return to StateCommunities, got %v", updatedModel.State)
	}
}
