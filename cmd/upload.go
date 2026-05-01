package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	targetID   string
	targetType int
)

var uploadCmd = &cobra.Command{
	Use:   "upload [file]",
	Short: "Upload media to a resource",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		if targetID == "" {
			return fmt.Errorf("target ID is required")
		}

		err := client.UploadMedia(targetID, targetType, filePath)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully uploaded %s to target %s\n", filePath, targetID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&targetID, "target", "t", "", "Target resource ID")
	uploadCmd.Flags().IntVarP(&targetType, "type", "y", 2, "Target resource type (1: Community, 2: Post, 3: Comment)")
}
