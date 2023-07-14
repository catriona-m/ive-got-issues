package app

import (
	"fmt"

	"github.com/google/go-github/v52/github"
	"github.com/ivegotissues/lib/gh"
)

type AddComment struct {
	LabelsFilter []string
	State        string
	Issues       []int
	Owner        string
	Comment      string
	Token        string
	Repo         string
	DryRun       bool
}

func (ac AddComment) AddComment() error {

	repo := gh.NewRepo(ac.Owner, ac.Repo, ac.Token)

	if len(ac.Issues) > 0 {
		err := ac.addCommentToIssueList(repo)
		if err != nil {
			return err
		}

	} else if len(ac.LabelsFilter) > 0 {
		err := ac.addCommentToIssuesFilteredByLabels(repo)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ac AddComment) addCommentToIssuesFilteredByLabels(repo gh.Repo) error {

	opts := github.IssueListByRepoOptions{
		State:  ac.State,
		Labels: ac.LabelsFilter,
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		issues, nextPage, err := repo.ListIssuesByRepo(opts)
		if err != nil {
			return fmt.Errorf("retrieving issues from github from page %d: %v", opts.ListOptions.Page, err)
		}
		for _, issue := range issues {

			fmt.Printf("Adding comment to issue: %d\t%s\t%s\n", issue.GetNumber(), issue.GetHTMLURL(), issue.GetTitle())

			if !ac.DryRun {
				err = repo.AddCommentToIssue(ac.Comment, issue.GetNumber())
				if err != nil {
					return err
				}
			}
		}

		if nextPage == 0 {
			break
		}
		opts.ListOptions.Page = nextPage

	}
	return nil
}

func (ac AddComment) addCommentToIssueList(repo gh.Repo) error {

	for _, issue := range ac.Issues {

		fmt.Printf("Adding comment to issue: %d\n", issue)

		if !ac.DryRun {
			err := repo.AddCommentToIssue(ac.Comment, issue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
