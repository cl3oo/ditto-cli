package ui

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcku/ditto/cli/internal/config"
)

func TestMainModel_CommandTransitions(t *testing.T) {
	cfg := &config.Config{BaseURL: "http://localhost:9001"}
	m := NewMainModel(cfg)

	tests := []struct {
		name     string
		command  string
		expected State
	}{
		{"Login command", ":login", StateLogin},
		{"Communities command", ":communities", StateLoading},
		{"Settings command", ":settings", StateProfileSettings},
		{"New post command", ":new", StateCreatePost},
		{"Following command", ":following", StateLoading},
		{"Joined command", ":joined", StateLoading},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.CommandBuffer = tt.command
			msg := tea.KeyMsg{Type: tea.KeyEnter}

			newModel, _ := m.Update(msg)
			updatedM := newModel.(MainModel)

			if updatedM.State != tt.expected {
				t.Errorf("%s: expected state %v, got %v", tt.name, tt.expected, updatedM.State)
			}
		})
	}

	t.Run("Post detail specific commands", func(t *testing.T) {
		m.State = StatePostDetail
		m.PostDetailModel.Post.ID = "p1"

		cmds := []struct {
			cmd      string
			expected State
		}{
			{":report spam", StateLoading},
			{":award silver", StateLoading},
			{":delete", StateLoading},
		}

		for _, tt := range cmds {
			m.CommandBuffer = tt.cmd
			newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
			updatedM := newModel.(MainModel)
			if updatedM.State != tt.expected {
				t.Errorf("%s: expected state %v, got %v", tt.cmd, tt.expected, updatedM.State)
			}
		}
	})
}

func TestMainModel_RefreshKeepsCurrentContext(t *testing.T) {
	tests := []struct {
		name          string
		state         State
		postID        string
		expectedPaths []string
	}{
		{
			name:          "feed refresh hits posts endpoint",
			state:         StateFeed,
			expectedPaths: []string{"/posts/"},
		},
		{
			name:          "communities refresh hits communities endpoint",
			state:         StateCommunities,
			expectedPaths: []string{"/communities/"},
		},
		{
			name:          "post detail refresh keeps current post detail",
			state:         StatePostDetail,
			postID:        "post-123",
			expectedPaths: []string{"/posts/post-123/", "/comments?id=post-123&type=2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requests []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path := r.URL.Path
				if r.URL.RawQuery != "" {
					path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
				}
				requests = append(requests, path)
				w.Header().Set("Content-Type", "application/json")
				switch path {
				case "/posts/", "/communities/":
					_, _ = w.Write([]byte(`{"data":[]}`))
				case "/posts/post-123/":
					_, _ = w.Write([]byte(`{"id":"post-123"}`))
				case "/comments?id=post-123&type=2":
					_, _ = w.Write([]byte(`{"data":[]}`))
				default:
					t.Fatalf("unexpected request path: %s", path)
				}
			}))
			defer server.Close()

			cfg := &config.Config{BaseURL: server.URL}
			m := NewMainModel(cfg)
			m.State = tt.state
			m.PostDetailModel.Post.ID = tt.postID
			m.CommandBuffer = ":refresh"

			newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
			updatedM := newModel.(MainModel)
			if updatedM.State != StateLoading {
				t.Fatalf("expected loading state, got %v", updatedM.State)
			}
			if cmd == nil {
				t.Fatal("expected refresh command")
			}

			_ = cmd()

			if len(requests) != len(tt.expectedPaths) {
				t.Fatalf("expected %d requests, got %d (%v)", len(tt.expectedPaths), len(requests), requests)
			}
			for i, expected := range tt.expectedPaths {
				if requests[i] != expected {
					t.Fatalf("expected request %d to be %s, got %s", i, expected, requests[i])
				}
			}
		})
	}
}
