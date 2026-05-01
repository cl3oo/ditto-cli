package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var outputPath string

var downloadCmd = &cobra.Command{
	Use:   "download [mediaID]",
	Short: "Download media from the platform",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mediaID := args[0]
		if outputPath == "" {
			outputPath = mediaID // Default to mediaID as filename
		}

		err := client.DownloadMedia(mediaID, outputPath)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully downloaded media to %s\n", outputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path")
}
