package cmd

import (
	"github.com/Songmu/prompter"
	c "github.com/gookit/color"
	"github.com/ive-got-issues/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// labelCmd represents the addLabels command
var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "label issues",
	Long:  `Adds specified labels to all issues in a repo based on input criteria such as state and content of the issue.`,
	Run: func(cmd *cobra.Command, args []string) {

		c.Info.Printf("Running label...\n")

		// TODO allow everything to be set via env var or config file too
		// TODO validate - can cobra do this for us?
		labels, _ := cmd.Flags().GetStringSlice("labels")
		content, _ := cmd.Flags().GetString("content")
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
			c.Info.Println("This is a dry-run only - to actually add labels on the issues please use --dry-run=false")
		} else {
			c.Warn.Println("This is NOT a dry-run - actual labels will be added to issues")
		}

		if openIssues && batch == 0 {
			c.Warn.Println("A browser tab will be opened for each issue without prompting. Use --batch to only open a specified number at a time")
			continueWithoutBatch := prompter.YN("Do you want to continue without --batch?", false)
			if !continueWithoutBatch {
				os.Exit(0)
			}
		}

		al := cli.AddLabels{
			Labels:     labels,
			Content:    content,
			State:      state,
			Issues:     issues,
			Batch:      batch,
			OpenIssues: openIssues,
			Token:      token,
			Repo:       repo,
			DryRun:     dryRun,
		}

		err := al.AddLabelsToIssues()
		if err != nil {
			c.Error.Printf("running label: %v", err)
			os.Exit(1)
		}

	},
}

func init() {

	rootCmd.AddCommand(labelCmd)
	labelCmd.Flags().StringSliceP("labels", "l", []string{}, "Labels to add")
	labelCmd.Flags().StringP("content", "m", "", "Filter which issues to add labels to based on whether the content body matches input string. Cannot be used if using issues flag")
	labelCmd.Flags().StringP("state", "s", "", "Filter which issues based on state. Possible values are 'open', 'closed' and 'all'")
	labelCmd.Flags().IntSliceP("issues", "i", []int{}, "List of issue numbers to add labels to. Cannot be used if using content")
	labelCmd.Flags().BoolP("dry-run", "d", true, "Print to console a simulation of what is expected to happen without making any actual changes to the issues. Defaults to true.")
	labelCmd.Flags().BoolP("open", "", false, "Open a browser tab for each issue labeled. Defaults to false.")
	labelCmd.Flags().IntP("batch", "b", 0, "Specify a number of issues to label at a time. If set to 0, all issues are labeled in one go. This setting is recommended when using open. Defaults to 0.")
	labelCmd.Flags().StringP("gh-repo", "", "", "The name of the github repo in the format 'owner/repo-name'")

	labelCmd.MarkFlagRequired("labels")
	labelCmd.MarkFlagRequired("gh-owner")
	labelCmd.MarkFlagRequired("gh-repo")
	labelCmd.MarkFlagsMutuallyExclusive("content", "issues")

}
