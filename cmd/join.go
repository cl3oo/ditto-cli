package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join a community",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := client.JoinCommunity(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Joined community %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
