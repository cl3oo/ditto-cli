package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {
	// Create a temporary home directory for testing
	tmpHome, err := os.MkdirTemp("", "ditto-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpHome)

	// Mock home directory by setting HOME environment variable
	// Note: On Windows this might need to be USERPROFILE, but os.UserHomeDir handles it
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	// Test default config
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.BaseURL != DefaultBaseURL {
		t.Errorf("Expected default BaseURL %s, got %s", DefaultBaseURL, cfg.BaseURL)
	}

	// Test saving and loading
	cfg.Token = "test-token"
	cfg.BaseURL = "http://test-api.local"
	err = cfg.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	cfg2, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig after save failed: %v", err)
	}
	if cfg2.Token != "test-token" {
		t.Errorf("Expected token test-token, got %s", cfg2.Token)
	}
	if cfg2.BaseURL != "http://test-api.local" {
		t.Errorf("Expected BaseURL http://test-api.local, got %s", cfg2.BaseURL)
	}

	// Test UpdateToken
	err = cfg2.UpdateToken("new-token")
	if err != nil {
		t.Fatalf("UpdateToken failed: %v", err)
	}
	cfg3, _ := LoadConfig()
	if cfg3.Token != "new-token" {
		t.Errorf("Expected new-token, got %s", cfg3.Token)
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}
	if filepath.Base(path) != "config.json" {
		t.Errorf("Expected config.json, got %s", filepath.Base(path))
	}
}
