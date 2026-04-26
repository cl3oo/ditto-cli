package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	voteValue int
)

var voteCmd = &cobra.Command{
	Use:   "vote",
	Short: "Vote on a resource",
}

var votePostCmd = &cobra.Command{
	Use:   "post [id]",
	Short: "Vote on a post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := client.Vote(args[0], 0, voteValue) // 0 = post
		if err != nil {
			return err
		}
		fmt.Printf("Voted %d on post %s\n", voteValue, args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(voteCmd)
	voteCmd.AddCommand(votePostCmd)
	voteCmd.PersistentFlags().IntVar(&voteValue, "value", 1, "Vote value (1 for upvote, -1 for downvote)")
}
