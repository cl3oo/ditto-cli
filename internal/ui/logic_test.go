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
				t.Errorf("expected state %v, got %v", tt.expected, updatedM.State)
			}
		})
	}
}
