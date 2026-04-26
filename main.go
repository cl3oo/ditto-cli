package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcku/ditto/cli/internal/config"
	"github.com/rfcku/ditto/cli/internal/ui"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v", err)
		os.Exit(1)
	}

	baseURL := os.Getenv("DITTO_API_URL")
	if baseURL == "" {
		if cfg.BaseURL != "" {
			baseURL = cfg.BaseURL
		} else {
			baseURL = "http://localhost:9001"
		}
	}

	m := ui.NewMainModel(baseURL)
	if cfg.Token != "" {
		m.Client.SetToken(cfg.Token)
		m.State = ui.StateLoading // Start by loading feed
	} else {
		m.State = ui.StateLogin
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// Save token on exit if it changed
	if m, ok := finalModel.(ui.MainModel); ok {
		if m.Client.Token != cfg.Token {
			cfg.Token = m.Client.Token
			cfg.BaseURL = m.Client.BaseURL
			if err := cfg.Save(); err != nil {
				fmt.Printf("Error saving config: %v", err)
			}
		}
	}
}
