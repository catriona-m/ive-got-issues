package cmd

import (
	"github.com/Songmu/prompter"
	c "github.com/gookit/color"
	"github.com/ive-got-issues/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var commentCmd = &cobra.Command{
	Use: "comment",
	// TODO descriptions
	Short: "Comment on issues",
	Long:  `Add a comment to all issues in a specified repo that match input filters such as state and labels.`,
	Run: func(cmd *cobra.Command, args []string) {
		c.Info.Printf("Running comment...\n")

		// TODO allow everything to be set via env var or config file too
		// TODO validate - can cobra do this for us?
		labels, _ := cmd.Flags().GetStringSlice("labels")
		comment, _ := cmd.Flags().GetString("comment")
		state, _ := cmd.Flags().GetString("state")
		issues, _ := cmd.Flags().GetIntSlice("issues")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		openIssues, _ := cmd.Flags().GetBool("open")
		batch, _ := cmd.Flags().GetInt("batch")
		repo, _ := cmd.Flags().GetString("gh-repo")

		// env vars
		viper.AutomaticEnv()
		token := viper.GetString("IGI_GITHUB_TOKEN")
		if token == "" {
			c.Error.Printf("Missing required environment variable `IGI_GITHUB_TOKEN`")
			os.Exit(1)
		}

		if dryRun {
			c.Info.Println("This is a dry-run only - to make actual comments on issues please use --dry-run=false")
		} else {
			c.Warn.Println("This is NOT a dry-run - actual comments will be added to issues")
		}

		if openIssues && batch == 0 {
			c.Warn.Println("A browser tab will be opened for each issue without prompting. Use --batch to only open a specified number at a time")
			continueWithoutBatch := prompter.YN("Do you want to continue without --batch?", false)
			if !continueWithoutBatch {
				os.Exit(0)
			}
		}

		ac := cli.AddComment{
			Labels:     labels,
			State:      state,
			Issues:     issues,
			Batch:      batch,
			OpenIssues: openIssues,
			Comment:    comment,
			Token:      token,
			Repo:       repo,
			DryRun:     dryRun,
		}

		err := ac.AddComment()
		if err != nil {
			c.Error.Printf("running comment: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)

	commentCmd.Flags().StringSliceP("labels", "l", []string{}, "Filters issues to comment on if they contain the specified labels. Cannot be used if using issues flag")
	commentCmd.Flags().StringP("state", "s", "", "Filter which issues based on state. Possible values are 'open', 'closed' and 'all'")
	commentCmd.Flags().IntSliceP("issues", "i", []int{}, "List of issue numbers to add labels to. Cannot be used if using labels flag")
	commentCmd.Flags().StringP("comment", "c", "", "Comment to add to issues")
	commentCmd.Flags().BoolP("dry-run", "d", true, "Print to console a simulation of what is expected to happen without making any actual changes to the issues. Defaults to true.")
	commentCmd.Flags().BoolP("open", "", false, "Open a browser tab for each issue commented. Defaults to false.")
	commentCmd.Flags().IntP("batch", "b", 0, "Specify a number of issues to comment at a time. If set to 0, all issues are commented in one go. This setting is recommended when using open. Defaults to 0.")
	commentCmd.Flags().StringP("gh-repo", "", "", "The name of the github repo in the format 'owner/repo-name'")

	commentCmd.MarkFlagRequired("state")
	commentCmd.MarkFlagRequired("gh-repo")
	commentCmd.MarkFlagRequired("comment")
	commentCmd.MarkFlagsMutuallyExclusive("labels", "issues")

}
