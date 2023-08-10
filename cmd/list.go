package cmd

import (
	"github.com/Songmu/prompter"
	c "github.com/gookit/color"
	"github.com/ive-got-issues/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists github issues",
	Long:  `Lists github issues based on labels, state and linked pull requests.`,
	Run: func(cmd *cobra.Command, args []string) {

		c.Info.Printf("Running list...\n")

		labels, _ := cmd.Flags().GetStringSlice("labels")
		content, _ := cmd.Flags().GetString("content")
		issueState, _ := cmd.Flags().GetString("state")
		prState, _ := cmd.Flags().GetString("pr-state")
		linkedPRs, _ := cmd.Flags().GetBool("linked-prs")
		openIssues, _ := cmd.Flags().GetBool("open")
		openPrs, _ := cmd.Flags().GetBool("open-prs")
		batch, _ := cmd.Flags().GetInt("batch")
		repo, _ := cmd.Flags().GetString("gh-repo")

		// env vars
		viper.AutomaticEnv()
		token := viper.GetString("IGI_GITHUB_TOKEN")
		if token == "" {
			c.Error.Printf("Missing required environment variable `IGI_GITHUB_TOKEN`")
			os.Exit(1)
		}

		if openIssues && batch == 0 {
			c.Warn.Println("A browser tab will be opened for each issue without prompting. Use --batch to only open a specified number at a time")
			continueWithoutBatch := prompter.YN("Do you want to continue without --batch?", false)
			if !continueWithoutBatch {
				os.Exit(0)
			}
		}

		li := cli.ListIssues{
			Labels:     labels,
			Content:    content,
			OpenIssues: openIssues,
			OpenPrs:    openPrs,
			Batch:      batch,
			LinkedPrs:  linkedPRs,
			PrState:    prState,
			IssueState: issueState,
			Token:      token,
			Repo:       repo,
		}

		err := li.ListIssues()
		if err != nil {
			c.Error.Printf("running list: %v", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceP("labels", "l", []string{}, "Labels to filter issues on")
	listCmd.Flags().StringP("content", "c", "", "Filter which issues to list based on whether the content body matches input string.")
	listCmd.Flags().StringP("state", "", "", "Filter which issues to list based on issue state. Possible values are 'open', 'closed' and 'all'")
	listCmd.Flags().BoolP("linked-prs", "p", false, "List matching issues that have linked Pull Requests associated with them. Defaults to false.")
	listCmd.Flags().StringP("pr-state", "", "", "Filter which issues and linked prs are listed based on the state of linked prs. Can only be used when linked-prs is 'true'. Possible values are 'open', 'closed' and 'merged'. If not specified, all linked pull requests are listed.")
	listCmd.Flags().BoolP("open", "", false, "Open a browser tab for each issue found. Defaults to false.")
	listCmd.Flags().BoolP("open-prs", "", false, "Open a browser tab for each linked PR found. Must be used in conjunction with 'linked-prs'. Defaults to false.")
	listCmd.Flags().IntP("batch", "b", 0, "Specify a number of issues to list at a time. If set to 0, all issues are listed. This setting is recommended when using open/open-prs. Defaults to 0.")
	listCmd.Flags().StringP("gh-repo", "", "", "The name of the github repo in the format 'owner/repo-name'")

	listCmd.MarkFlagRequired("state")
	listCmd.MarkFlagRequired("gh-repo")

}
