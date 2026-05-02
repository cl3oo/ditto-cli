package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func maskToken(token string) string {
	if len(token) <= 8 {
		return token
	}
	return token[:8] + "..."
}

var (
	username string
	password string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Ditto",
	RunE: func(cmd *cobra.Command, args []string) error {
		if username == "" {
			return fmt.Errorf("username is required")
		}

		resolvedPassword, err := resolvePassword(cmd.OutOrStdout())
		if err != nil {
			return err
		}

		token, err := client.Login(username, resolvedPassword)
		if err != nil {
			return err
		}

		cfg.Token = token
		if baseURL != "" {
			cfg.BaseURL = baseURL
		}
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}

		fmt.Printf("Logged in successfully. Token: %s\n", maskToken(token))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
}
