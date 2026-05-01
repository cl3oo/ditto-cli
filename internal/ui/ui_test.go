package ui

import (
	"strings"
	"testing"

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
