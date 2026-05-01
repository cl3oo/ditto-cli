package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/rfcku/ditto-cli/internal/api"
	"github.com/rfcku/ditto-cli/internal/config"
)

var (
	cfg     *config.Config
	client  *api.Client
	baseURL string
	token   string
)

var rootCmd = &cobra.Command{
	Use:   "ditto",
	Short: "Ditto CLI - A command line interface for Ditto",
	Long:  `Ditto CLI allows you to interact with the Ditto API and launch the TUI.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		if baseURL == "" {
			baseURL = os.Getenv("DITTO_API_URL")
			if baseURL == "" {
				baseURL = cfg.BaseURL
			}
		}

		if token == "" {
			token = cfg.Token
		}

		client = api.NewClient(baseURL)
		if token != "" {
			client.SetToken(token)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// If no command is provided, we can either show help or run TUI.
		// User said "if the user executes 'ditto tui' it will launch the tui app", 
		// so I'll show help by default if no args.
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&baseURL, "url", "", "API base URL")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "API token")
}
