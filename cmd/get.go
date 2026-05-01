package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many resources",
	Long:  `Get retrieves resources from the Ditto API.`,
}

var getPostsCmd = &cobra.Command{
	Use:   "posts",
	Short: "Get posts from trending feed",
	RunE: func(cmd *cobra.Command, args []string) error {
		posts, err := client.GetTrendingPosts()
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tAUTHOR\tCOMMUNITY\tSCORE")
		for _, p := range posts {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
				p.ID, p.Title, p.Author.Username, p.Community.Name, p.Scores.VoteScore)
		}
		return w.Flush()
	},
}

var getFeedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Get posts from personalized feed",
	RunE: func(cmd *cobra.Command, args []string) error {
		posts, err := client.GetFeed()
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tAUTHOR\tCOMMUNITY\tSCORE")
		for _, p := range posts {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
				p.ID, p.Title, p.Author.Username, p.Community.Name, p.Scores.VoteScore)
		}
		return w.Flush()
	},
}

var getCommunitiesCmd = &cobra.Command{

	Use:   "communities",
	Short: "List all communities",
	RunE: func(cmd *cobra.Command, args []string) error {
		communities, err := client.GetCommunities()
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTITLE\tMEMBERS")
		for _, c := range communities {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", 
				c.ID, c.Name, c.Title, c.Scores.SubCount)
		}
		return w.Flush()
	},
}

var getMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Get current user profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		user, err := client.GetMe()
		if err != nil {
			return err
		}

		data, _ := json.MarshalIndent(user, "", "  ")
		fmt.Println(string(data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getPostsCmd)
	getCmd.AddCommand(getFeedCmd)
	getCmd.AddCommand(getCommunitiesCmd)
	getCmd.AddCommand(getMeCmd)
}

