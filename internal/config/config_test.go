package config

import (
	"os"
	"path/filepath"
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
	tmpFile := filepath.Join(t.TempDir(), "config_test.toml")

	ConfigPathOverride = tmpFile
	defer func() { ConfigPathOverride = "" }()

	cfg := DefaultConfig()
	cfg.Token = "test-token"
	cfg.BaseURL = "http://test-api"

	err := cfg.Save()
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loaded.Token != "test-token" {
		t.Errorf("expected token test-token, got %s", loaded.Token)
	}
	if loaded.BaseURL != "http://test-api" {
		t.Errorf("expected base URL http://test-api, got %s", loaded.BaseURL)
	}

	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("failed to stat config: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected config file mode 0600, got %o", info.Mode().Perm())
	}
}
