package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rfcku/ditto/cli/internal/ui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the TUI application",
	Run: func(cmd *cobra.Command, args []string) {
		m := ui.NewMainModel(baseURL)
		if token != "" {
			m.Client.SetToken(token)
			m.State = ui.StateLoading
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
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
