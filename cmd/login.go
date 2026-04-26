package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	username string
	password string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Ditto",
	RunE: func(cmd *cobra.Command, args []string) error {
		if username == "" || password == "" {
			return fmt.Errorf("username and password are required")
		}

		token, err := client.Login(username, password)
		if err != nil {
			return err
		}

		cfg.Token = token
		cfg.BaseURL = baseURL
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}

		fmt.Println("Logged in successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
}
