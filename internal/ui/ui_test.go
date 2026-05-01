package ui

import (
	"testing"

	"github.com/rfcku/ditto/cli/internal/api"
	"github.com/rfcku/ditto/cli/internal/config"
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
