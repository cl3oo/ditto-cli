package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Show details of a specific resource",
}

var describePostCmd = &cobra.Command{
	Use:   "post [id]",
	Short: "Describe a post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		post, err := client.GetPost(args[0])
		if err != nil {
			return err
		}

		data, _ := json.MarshalIndent(post, "", "  ")
		fmt.Println(string(data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.AddCommand(describePostCmd)
}
