package cli

import (
	"fmt"
	"regexp"

	"github.com/Songmu/prompter"
	"github.com/google/go-github/v52/github"
	c "github.com/gookit/color"
	"github.com/ive-got-issues/lib/gh"
	"github.com/pkg/browser"
)

type AddLabels struct {
	Labels     []string
	Content    string
	State      string
	Issues     []int
	OpenIssues bool
	Batch      int
	Token      string
	Repo       string
	DryRun     bool
}

const githubUrl = "https://github.com/"

func (al AddLabels) AddLabelsToIssues() error {

	repo := gh.NewRepo(al.Repo, al.Token)

	if len(al.Issues) > 0 {
		counter, err := al.labelIssues(repo)
		if err != nil {
			return err
		}
		c.Info.Printf("Finished labelling %d issues", counter)
	} else if al.Content != "" {

		c.Info.Printf("Finding matching issues to label in %s", al.Repo)
		opts := github.IssueListByRepoOptions{
			State: al.State,
			ListOptions: github.ListOptions{
				PerPage: 100,
				Page:    1,
			},
		}

		issueCount := 0
		totalIssuesCount := 0
		batchCounter := 1
		var continueLabelling bool
		for {

			issues, nextPage, err := repo.ListIssuesByRepo(opts)
			if err != nil {
				return fmt.Errorf("retrieving issues from github from page %d: %v", opts.ListOptions.Page, err)
			}
			continueLabelling, issueCount, err = al.labelIssuesFilteredByBodyContent(repo, issues, batchCounter)
			if err != nil {
				return err
			}

			totalIssuesCount += issueCount

			if al.Batch > 0 && totalIssuesCount+1 >= al.Batch {
				batchCounter = (totalIssuesCount + 1) % al.Batch
				if batchCounter == 0 {
					batchCounter = al.Batch
				}
			} else {
				batchCounter = totalIssuesCount + 1
			}
			if !continueLabelling || nextPage == 0 {
				c.Info.Printf("Finished labelling %d issues", totalIssuesCount)
				break
			}
			opts.ListOptions.Page = nextPage
		}
	}

	return nil

}

func (al AddLabels) labelIssuesFilteredByBodyContent(repo gh.Repo, issues []*github.Issue, batchCounter int) (bool, int, error) {

	re, err := regexp.Compile(al.Content)
	if err != nil {
		return false, 0, fmt.Errorf("regex %s failed to compile: %v", al.Content, err)
	}

	issueCounter := 0
	for _, issue := range issues {
		if issue.IsPullRequest() {
			continue
		}

		if re.MatchString(issue.GetBody()) {

			// skip this one if the issue already has the labels added
			if !hasLabels(issue.Labels, al.Labels) {
				c.Printf("Adding labels <magenta>%v</> to issue: <cyan>%d</>\t%s\t%s\n", al.Labels, issue.GetNumber(), issue.GetHTMLURL(), issue.GetTitle())
				if !al.DryRun {
					err = repo.AddLabelsToIssue(al.Labels, issue.GetNumber())
					if err != nil {
						return false, issueCounter, nil
					}
				}

				issueCounter++

				if al.OpenIssues {
					if err := browser.OpenURL(issue.GetHTMLURL()); err != nil {
						c.Error.Printf("failed to open issue %s in browser: %v", issue.GetHTMLURL(), err)
					}
				}

				if al.Batch > 0 {
					if batchCounter == al.Batch {
						continueListing := prompter.YN("Do you want to continue labelling issues?", true)
						if !continueListing {
							return false, issueCounter, nil
						}
						batchCounter = 1
						continue
					}
					batchCounter++
				}
			}
		}
	}
	return true, issueCounter, nil
}

func (al AddLabels) labelIssues(repo gh.Repo) (int, error) {
	issueCounter := 0
	for _, issue := range al.Issues {
		c.Printf("Adding labels <magenta>%v</> to issue: <cyan>#%d</>\n", al.Labels, issue)
		if !al.DryRun {
			err := repo.AddLabelsToIssue(al.Labels, issue)
			if err != nil {
				return issueCounter, err
			}
		}

		issueCounter++

		if al.OpenIssues {
			url := fmt.Sprintf("%s%s/%s/issues/%d", githubUrl, repo.Owner, repo.Name, issue)
			if err := browser.OpenURL(url); err != nil {
				c.Error.Printf("failed to open issue %s in browser", url)
			}
		}

		if al.Batch > 0 {
			if issueCounter >= al.Batch {
				if issueCounter%al.Batch == 0 {
					msg := c.Sprintf("<green>Do you want to continue labelling issues?</>")
					continueListing := prompter.YN(msg, true)
					if !continueListing {
						return issueCounter, nil
					}
					continue
				}
			}
		}
	}
	return issueCounter, nil
}

func hasLabels(issueLabels []*github.Label, labelsToAdd []string) bool {
	for _, label := range labelsToAdd {
		hasLabel := false
		for _, issueLabel := range issueLabels {
			if issueLabel.GetName() == label {
				hasLabel = true
			}
		}
		if !hasLabel {
			return false
		}
	}
	return true
}
