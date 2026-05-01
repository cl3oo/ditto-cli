package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	email string
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if username == "" || password == "" {
			return fmt.Errorf("username and password are required")
		}

		if email == "" {
			email = username + "@example.com"
		}

		token, err := client.Register(username, email, password)
		if err != nil {
			return err
		}

		cfg.Token = token
		cfg.BaseURL = baseURL
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}

		fmt.Println("Account created and logged in successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVarP(&username, "username", "u", "", "Username")
	registerCmd.Flags().StringVarP(&email, "email", "e", "", "Email")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
}
