package cli

import (
	"fmt"

	"github.com/Songmu/prompter"
	"github.com/google/go-github/v52/github"
	c "github.com/gookit/color"
	"github.com/ive-got-issues/lib/gh"
	"github.com/pkg/browser"
)

type AddComment struct {
	Labels     []string
	State      string
	Issues     []int
	Comment    string
	Batch      int
	OpenIssues bool
	Token      string
	Repo       string
	DryRun     bool
}

func (ac AddComment) AddComment() error {

	repo := gh.NewRepo(ac.Repo, ac.Token)

	var counter int
	var err error
	if len(ac.Issues) > 0 {
		counter, err = ac.addCommentToIssueList(repo)
		if err != nil {
			return err
		}

	} else if len(ac.Labels) > 0 {
		c.Info.Printf("Finding matching issues to comment in %s\n", ac.Repo)
		counter, err = ac.addCommentToIssuesFilteredByLabels(repo)
		if err != nil {
			return err
		}
	}

	c.Info.Printf("Finished commenting %d issues", counter)

	return nil
}

func (ac AddComment) addCommentToIssuesFilteredByLabels(repo gh.Repo) (int, error) {

	opts := github.IssueListByRepoOptions{
		State:  ac.State,
		Labels: ac.Labels,
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	issueCounter := 0
	batchCounter := 1
	for {
		issues, nextPage, err := repo.ListIssuesByRepo(opts)
		if err != nil {
			return issueCounter, fmt.Errorf("retrieving issues from github from page %d: %v", opts.ListOptions.Page, err)
		}
		for _, issue := range issues {

			c.Printf("Adding comment to <cyan>#%d</>\t%s\t%s\n", issue.GetNumber(), issue.GetTitle(), issue.GetHTMLURL())

			if !ac.DryRun {
				err = repo.AddCommentToIssue(ac.Comment, issue.GetNumber())
				if err != nil {
					return issueCounter, err
				}
			}
			issueCounter++

			if ac.OpenIssues {
				if err := browser.OpenURL(issue.GetHTMLURL()); err != nil {
					c.Error.Printf("failed to open issue %s in browser: %v", issue.GetHTMLURL(), err)
				}
			}

			if ac.Batch > 0 {
				if batchCounter == ac.Batch {
					continueListing := prompter.YN("Do you want to continue commenting issues?", true)
					if !continueListing {
						return issueCounter, nil
					}
					batchCounter = 1
					continue
				}
				batchCounter++
			}
		}

		if nextPage == 0 {
			break
		}
		opts.ListOptions.Page = nextPage

	}
	return issueCounter, nil
}

func (ac AddComment) addCommentToIssueList(repo gh.Repo) (int, error) {
	count := 0
	for _, issue := range ac.Issues {

		c.Printf("Adding comment to <cyan>#%d</>\n", issue)

		if !ac.DryRun {
			err := repo.AddCommentToIssue(ac.Comment, issue)
			if err != nil {
				return count, err
			}
		}
		count++
	}
	return count, nil
}
