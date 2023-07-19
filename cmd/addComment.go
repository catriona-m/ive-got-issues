package cmd

import (
	"fmt"
	"github.com/ivegotissues/cli"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var addCommentCmd = &cobra.Command{
	Use: "addComment",
	// TODO descriptions
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running addComment...")

		// TODO allow everything to be set via env var or config file too
		// TODO validate - can cobra do this for us?
		labels, _ := cmd.Flags().GetStringSlice("labels")
		comment, _ := cmd.Flags().GetString("comment")
		state, _ := cmd.Flags().GetString("state")
		issues, _ := cmd.Flags().GetIntSlice("issues")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		openIssues, _ := cmd.Flags().GetBool("open-issues")
		batch, _ := cmd.Flags().GetInt("batch")
		owner, _ := cmd.Flags().GetString("gh-owner")
		repo, _ := cmd.Flags().GetString("gh-repo")

		// env vars
		viper.AutomaticEnv()
		token := viper.GetString("IGI_GITHUB_TOKEN")

		ac := cli.AddComment{
			Labels:     labels,
			State:      state,
			Issues:     issues,
			Batch:      batch,
			OpenIssues: openIssues,
			Owner:      owner,
			Comment:    comment,
			Token:      token,
			Repo:       repo,
			DryRun:     dryRun,
		}

		err := ac.AddComment()
		if err != nil {
			fmt.Errorf("running addComment: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCommentCmd)

	addCommentCmd.Flags().StringSliceP("labels", "l", []string{}, "Filters issues to comment on if they contain the specified labels. Cannot be used if using issues flag")
	addCommentCmd.Flags().StringP("state", "s", "", "Filter which issues based on state. Possible values are 'open', 'closed' and 'all'")
	addCommentCmd.Flags().IntSliceP("issues", "i", []int{}, "List of issue numbers to add labels to. Cannot be used if using labels flag")
	addCommentCmd.Flags().StringP("comment", "c", "", "Comment to add to issues")
	addCommentCmd.Flags().BoolP("dry-run", "d", true, "Print to console a simulation of what is expected to happen without making any actual changes to the issues. Defaults to true.")
	addCommentCmd.Flags().BoolP("open-issues", "", false, "Open a browser tab for each issue commented. Defaults to false.")
	addCommentCmd.Flags().IntP("batch", "b", 0, "Specify a number of issues to comment at a time. If set to 0, all issues are commented in one go. This setting is recommended when using open-issues. Defaults to 0.")
	addCommentCmd.Flags().StringP("gh-owner", "", "", "The name of the github owner")
	addCommentCmd.Flags().StringP("gh-repo", "", "", "The name of the github repo")

	addCommentCmd.MarkFlagRequired("state")
	addCommentCmd.MarkFlagRequired("gh-owner")
	addCommentCmd.MarkFlagRequired("gh-repo")
	addCommentCmd.MarkFlagRequired("comment")
	addCommentCmd.MarkFlagsMutuallyExclusive("labels", "issues")

}
