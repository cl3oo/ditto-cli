package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcku/ditto-cli/internal/config"
	"github.com/rfcku/ditto-cli/internal/types"
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
			{":delete", StateConfirm},
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

	t.Run("Delete account requires double confirmation", func(t *testing.T) {
		m.State = StateProfileSettings
		m.Me = &types.User{ID: "u1", Username: "tester"}
		m.CommandBuffer = ":delete-account"

		newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		updatedM := newModel.(MainModel)
		if updatedM.State != StateConfirm {
			t.Fatalf("expected confirm state, got %v", updatedM.State)
		}
		if updatedM.ConfirmDialog.Step != 1 || updatedM.ConfirmDialog.Steps != 2 {
			t.Fatalf("expected double confirm step 1/2, got %d/%d", updatedM.ConfirmDialog.Step, updatedM.ConfirmDialog.Steps)
		}

		confirmedModel, _ := updatedM.Update(tea.KeyMsg{Type: tea.KeyEnter})
		confirmed := confirmedModel.(MainModel)
		if confirmed.State != StateConfirm {
			t.Fatalf("expected second confirm prompt, got %v", confirmed.State)
		}
		if confirmed.ConfirmDialog.Step != 2 {
			t.Fatalf("expected second confirm step, got %d", confirmed.ConfirmDialog.Step)
		}
	})
}
