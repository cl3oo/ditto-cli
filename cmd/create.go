package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	postTitle     string
	postContent   string
	postCommunity string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource",
}

var createPostCmd = &cobra.Command{
	Use:   "post",
	Short: "Create a new post",
	RunE: func(cmd *cobra.Command, args []string) error {
		if postTitle == "" || postContent == "" || postCommunity == "" {
			return fmt.Errorf("title, content and community are required")
		}

		err := client.CreatePost(postTitle, postContent, postCommunity)
		if err != nil {
			return err
		}

		fmt.Println("Post created successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createPostCmd)

	createPostCmd.Flags().StringVarP(&postTitle, "title", "t", "", "Post title")
	createPostCmd.Flags().StringVarP(&postContent, "content", "c", "", "Post content")
	createPostCmd.Flags().StringVarP(&postCommunity, "community", "m", "", "Community name")
}
