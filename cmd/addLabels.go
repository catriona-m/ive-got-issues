package cmd

import (
	"fmt"
	"github.com/ivegotissues/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addLabelsCmd represents the addLabels command
var addLabelsCmd = &cobra.Command{
	Use:   "addLabels",
	Short: "Adds labels to issues based on input criteria",
	// TODO long description
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Running addLabels...")

		// TODO allow everything to be set via env var or config file too
		// TODO validate - can cobra do this for us?
		labels, _ := cmd.Flags().GetStringSlice("labels")
		contentMatches, _ := cmd.Flags().GetString("content-matches")
		state, _ := cmd.Flags().GetString("state")
		issues, _ := cmd.Flags().GetIntSlice("issues")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		owner, _ := cmd.Flags().GetString("gh-owner")
		repo, _ := cmd.Flags().GetString("gh-repo")

		// env vars
		viper.AutomaticEnv()
		token := viper.GetString("IGI_GITHUB_TOKEN")

		al := app.AddLabels{
			Labels:         labels,
			ContentMatches: contentMatches,
			State:          state,
			Issues:         issues,
			Owner:          owner,
			Token:          token,
			Repo:           repo,
			DryRun:         dryRun,
		}

		err := al.AddLabelsToIssues()
		if err != nil {
			fmt.Errorf("running addLabels: %v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(addLabelsCmd)
	addLabelsCmd.Flags().StringSliceP("labels", "l", []string{}, "Labels to add")
	addLabelsCmd.Flags().StringP("content-matches", "m", "", "Filter which issues to add labels to based on whether the content body matches input string. Cannot be used if using issues flag")
	addLabelsCmd.Flags().StringP("state", "s", "", "Filter which issues based on state. Possible values are 'open', 'closed' and 'all'")
	addLabelsCmd.Flags().IntSliceP("issues", "i", []int{}, "List of issue numbers to add labels to. Cannot be used if using content-matches")
	addLabelsCmd.Flags().BoolP("dry-run", "d", true, "Print to console a simulation of what is expected to happen without making any actual changes to the issues. Defaults to true.")
	addLabelsCmd.Flags().StringP("gh-owner", "", "", "The name of the github owner")
	addLabelsCmd.Flags().StringP("gh-repo", "", "", "The name of the github repo")

	addLabelsCmd.MarkFlagRequired("labels")
	addLabelsCmd.MarkFlagRequired("state")
	addLabelsCmd.MarkFlagRequired("gh-owner")
	addLabelsCmd.MarkFlagRequired("gh-repo")
	addLabelsCmd.MarkFlagsMutuallyExclusive("content-matches", "issues")

}
