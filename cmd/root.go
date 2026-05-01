package cmd

import (
	"fmt"
	"os"

	"github.com/rfcku/ditto-cli/internal/api"
	"github.com/rfcku/ditto-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg       *config.Config
	client    *api.Client
	baseURL   string
	token     string
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "ditto-cli",
	Short:   "Ditto CLI - A command line interface for Ditto",
	Long:    `Ditto CLI allows you to interact with the Ditto API and launch the TUI.`,
	Version: versionString(),
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
		// so show help by default if no args.
		_ = cmd.Help()
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
	rootCmd.AddCommand(versionCmd)
}

func versionString() string {
	return fmt.Sprintf("%s (commit %s, built %s)", version, commit, buildDate)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build version information",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Fprintln(cmd.OutOrStdout(), versionString())
		return err
	},
}
