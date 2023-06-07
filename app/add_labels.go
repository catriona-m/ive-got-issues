package app

import (
	"fmt"
	"regexp"

	"github.com/google/go-github/v52/github"
	"github.com/ivegotissues/lib/gh"
)

type AddLabels struct {
	Labels         []string
	ContentMatches string
	State          string
	Issues         []int
	Owner          string
	Token          string
	Repo           string
	DryRun         bool
}

func (al AddLabels) AddLabelsToIssues() error {

	repo := gh.NewRepo(al.Owner, al.Repo, al.Token)

	if len(al.Issues) > 0 {
		err := al.labelIssues(repo)
		if err != nil {
			return err
		}
	} else if al.ContentMatches != "" {

		opts := github.IssueListByRepoOptions{
			State: al.State,
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
			err = al.labelIssuesFilteredByBodyContent(repo, issues)
			if err != nil {
				return err
			}

			if nextPage == 0 {
				break
			}
			opts.ListOptions.Page = nextPage
		}
	}

	return nil

}

func (al AddLabels) labelIssuesFilteredByBodyContent(repo gh.Repo, issues []*github.Issue) error {

	re, err := regexp.Compile(al.ContentMatches)
	if err != nil {
		return fmt.Errorf("regex %s failed to compile: %v", al.ContentMatches, err)
	}

	for _, issue := range issues {
		if issue.IsPullRequest() {
			continue
		}
		if issue.Body != nil {
			if re.MatchString(issue.GetBody()) {

				// skip this one if the issue already has the labels added
				if !hasLabels(issue.Labels, al.Labels) {
					fmt.Printf("Adding labels %v to issue: %d\t%s\t%s\n", al.Labels, issue.GetNumber(), issue.GetHTMLURL(), issue.GetTitle())
					if !al.DryRun {
						err = repo.AddLabelsToIssue(al.Labels, issue.GetNumber())
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func (al AddLabels) labelIssues(repo gh.Repo) error {
	for _, issue := range al.Issues {
		fmt.Printf("Adding labels %v to issue: %d\n", al.Labels, issue)
		if !al.DryRun {
			err := repo.AddLabelsToIssue(al.Labels, issue)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
