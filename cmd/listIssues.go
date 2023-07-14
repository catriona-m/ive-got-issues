package cmd

import (
	"fmt"
	"github.com/ivegotissues/app"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// listIssuesCmd represents the listIssues command
var listIssuesCmd = &cobra.Command{
	Use:   "listIssues",
	Short: "Lists github issues",
	Long:  `Lists github issues based on labels, state and linked pull requests.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running listIssues...")

		labels, _ := cmd.Flags().GetStringSlice("labels-filter")
		contentFilter, _ := cmd.Flags().GetString("content-filter")
		issueState, _ := cmd.Flags().GetString("issue-state")
		prState, _ := cmd.Flags().GetString("pr-state")
		linkedPRs, _ := cmd.Flags().GetBool("linked-prs")
		browseIssues, _ := cmd.Flags().GetBool("browse-issues")
		browsePrs, _ := cmd.Flags().GetBool("browse-prs")
		batch, _ := cmd.Flags().GetInt("batch")
		owner, _ := cmd.Flags().GetString("gh-owner")
		repo, _ := cmd.Flags().GetString("gh-repo")

		// env vars
		viper.AutomaticEnv()
		token := viper.GetString("IGI_GITHUB_TOKEN")

		li := app.ListIssues{
			LabelsFilter:  labels,
			ContentFilter: contentFilter,
			BrowseIssues:  browseIssues,
			BrowsePrs:     browsePrs,
			Batch:         batch,
			LinkedPrs:     linkedPRs,
			PrState:       prState,
			IssueState:    issueState,
			Owner:         owner,
			Token:         token,
			Repo:          repo,
		}

		err := li.ListIssues()
		if err != nil {
			fmt.Errorf("running listIssues: %v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(listIssuesCmd)

	listIssuesCmd.Flags().StringSliceP("labels-filter", "l", []string{}, "Labels to filter issues on")
	listIssuesCmd.Flags().StringP("content-filter", "c", "", "Filter which issues to list based on whether the content body matches input string.")
	listIssuesCmd.Flags().StringP("issue-state", "", "", "Filter which issues to list based on issue state. Possible values are 'open', 'closed' and 'all'")
	listIssuesCmd.Flags().BoolP("linked-prs", "p", false, "List matching issues that have linked Pull Requests associated with them. Defaults to false.")
	listIssuesCmd.Flags().StringP("pr-state", "", "", "Filter which issues and linked prs are listed based on the state of linked prs. Can only be used when linked-prs is 'true'. Possible values are 'open', 'closed' and 'merged'. If not specified, all linked pull requests are listed.")
	listIssuesCmd.Flags().BoolP("browse-issues", "", false, "Open a browser tab for each issue found. Defaults to false.")
	listIssuesCmd.Flags().BoolP("browse-prs", "", false, "Open a browser tab for each linked PR found. Must be used in conjunction with 'linked-prs'. Defaults to false.")
	listIssuesCmd.Flags().IntP("batch", "b", 0, "Specify a number of issues to list at a time. If set to 0, all issues are listed. This setting is recommended when using browse-issues/prs. Defaults to 0.")
	listIssuesCmd.Flags().StringP("gh-owner", "", "", "The name of the github owner")
	listIssuesCmd.Flags().StringP("gh-repo", "", "", "The name of the github repo")

	listIssuesCmd.MarkFlagRequired("issue-state")
	listIssuesCmd.MarkFlagRequired("gh-owner")
	listIssuesCmd.MarkFlagRequired("gh-repo")

}
