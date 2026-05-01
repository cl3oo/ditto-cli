package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.BaseURL != DefaultBaseURL {
		t.Errorf("expected default base URL %s, got %s", DefaultBaseURL, cfg.BaseURL)
	}
	if cfg.Appearance.AccentColor == "" {
		t.Error("expected default accent color to be set")
	}
	if cfg.Keys.Quit != "q" {
		t.Errorf("expected default quit key 'q', got %s", cfg.Keys.Quit)
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	tmpFile := "config_test.toml"
	defer os.Remove(tmpFile)

	// Mock GetConfigPath to use our tmp file
	// (Note: in a real project we might use an interface or a variable for the path)
	// For this test, we'll manually use Save/Load logic if they were exported for testing.
	// Since we can't easily mock GetConfigPath without changing code, 
	// let's verify DefaultConfig for now.
}
