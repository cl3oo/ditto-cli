package cmd

import (
	"testing"

	"github.com/rfcku/ditto-cli/internal/api"
	"github.com/rfcku/ditto-cli/internal/config"
)

func TestRegisterRequiresExplicitEmail(t *testing.T) {
	prevUsername, prevEmail, prevPassword := username, email, password
	prevCfg, prevClient, prevBaseURL := cfg, client, baseURL
	defer func() {
		username, email, password = prevUsername, prevEmail, prevPassword
		cfg, client, baseURL = prevCfg, prevClient, prevBaseURL
	}()

	cfg = config.DefaultConfig()
	client = api.NewClient("http://localhost:9001/v1")
	baseURL = ""
	username = "alice"
	email = ""
	password = "secret"

	err := registerCmd.RunE(registerCmd, nil)
	if err == nil {
		t.Fatal("expected missing email to fail")
	}
	if got := err.Error(); got != "username and email are required" {
		t.Fatalf("unexpected error: %s", got)
	}
}
