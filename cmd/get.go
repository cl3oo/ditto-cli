package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rfcku/ditto-cli/internal/types"
	"github.com/spf13/cobra"
)

var randomNum int

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many resources",
	Long:  `Get retrieves resources from the Ditto API.`,
}

var getPostsCmd = &cobra.Command{
	Use:   "posts",
	Short: "Get posts",
	RunE: func(cmd *cobra.Command, args []string) error {
		var posts []types.Post
		var err error
		if randomNum > 0 {
			posts, err = client.GetRandomPosts(randomNum)
		} else {
			posts, err = client.GetTrendingPosts()
		}
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "ID\tTITLE\tAUTHOR\tCOMMUNITY\tSCORE")
		for _, p := range posts {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
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
		_, _ = fmt.Fprintln(w, "ID\tTITLE\tAUTHOR\tCOMMUNITY\tSCORE")
		for _, p := range posts {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
				p.ID, p.Title, p.Author.Username, p.Community.Name, p.Scores.VoteScore)
		}
		return w.Flush()
	},
}

var getCommunitiesCmd = &cobra.Command{
	Use:   "communities",
	Short: "List all communities",
	RunE: func(cmd *cobra.Command, args []string) error {
		var communities []types.Community
		var err error
		if randomNum > 0 {
			communities, err = client.GetRandomCommunities(randomNum)
		} else {
			communities, err = client.GetCommunities()
		}
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		_, _ = fmt.Fprintln(w, "ID\tNAME\tTITLE\tMEMBERS")
		for _, c := range communities {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
				c.ID, c.Name, c.Title, c.Scores.SubCount)
		}
		return w.Flush()
	},
}

var getCommentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Get comments",
	RunE: func(cmd *cobra.Command, args []string) error {
		if randomNum > 0 {
			comments, err := client.GetRandomComments(randomNum)
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, "ID\tCONTENT\tAUTHOR\tSCORE")
			for _, c := range comments {
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
					c.ID, c.Content, c.Author.Username, c.Scores.VoteScore)
			}
			return w.Flush()
		}
		return fmt.Errorf("please use --random flag for now")
	},
}

var getUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Get users",
	RunE: func(cmd *cobra.Command, args []string) error {
		if randomNum > 0 {
			users, err := client.GetRandomUsers(randomNum)
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
			_, _ = fmt.Fprintln(w, "ID\tUSERNAME")
			for _, u := range users {
				_, _ = fmt.Fprintf(w, "%s\t%s\n", u.ID, u.Username)
			}
			return w.Flush()
		}
		return fmt.Errorf("please use --random flag for now")
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
	getCmd.PersistentFlags().IntVarP(&randomNum, "random", "r", 0, "Number of random items to fetch")
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getPostsCmd)
	getCmd.AddCommand(getFeedCmd)
	getCmd.AddCommand(getCommunitiesCmd)
	getCmd.AddCommand(getMeCmd)
	getCmd.AddCommand(getCommentsCmd)
	getCmd.AddCommand(getUsersCmd)
}
