package ui

import (
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
