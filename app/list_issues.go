package app

import (
	"fmt"
	"regexp"

	"github.com/Songmu/prompter"
	"github.com/google/go-github/v52/github"
	"github.com/ivegotissues/lib/gh"
	"github.com/pkg/browser"
)

type ListIssues struct {
	LabelsFilter  []string
	ContentFilter string
	IssueState    string
	LinkedPrs     bool
	PrState       string
	BrowseIssues  bool
	BrowsePrs     bool
	Batch         int
	Owner         string
	Token         string
	Repo          string
}

func (li ListIssues) ListIssues() error {

	opts := github.IssueListByRepoOptions{
		State:  li.IssueState,
		Labels: li.LabelsFilter,
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	repo := gh.NewRepo(li.Owner, li.Repo, li.Token)

	for {

		issues, nextPage, err := repo.ListIssuesByRepo(opts)
		if err != nil {
			return fmt.Errorf("retrieving issues from github from page %d: %v", opts.ListOptions.Page, err)
		}
		continueProcessing, err := li.listFilteredIssues(repo, issues)
		if err != nil {
			return err
		}

		if !continueProcessing || nextPage == 0 {
			fmt.Println("Finished listing issues.")
			break
		}
		opts.ListOptions.Page = nextPage
	}

	return nil
}

func (li ListIssues) listFilteredIssues(repo gh.Repo, issues []*github.Issue) (bool, error) {

	counter := 1
	var err error

	re := &regexp.Regexp{}
	if li.ContentFilter != "" {
		re, err = regexp.Compile(li.ContentFilter)
		if err != nil {
			return false, err
		}
	}

	for _, issue := range issues {
		if issue.IsPullRequest() {
			continue
		}

		if li.ContentFilter != "" {
			if !re.MatchString(issue.GetBody()) {
				continue
			}
		}

		prs := make(map[string]string)
		if li.LinkedPrs {
			prs, err = crossReferencedPRs(repo, issue.GetNumber(), li.PrState)
			if err != nil {
				return false, err
			}
			if len(prs) == 0 {
				continue
			}
		}

		fmt.Printf("Issue: %d\t%s\t%s\n", issue.GetNumber(), issue.GetHTMLURL(), issue.GetTitle())

		if li.BrowseIssues {
			if err := browser.OpenURL(issue.GetHTMLURL()); err != nil {
				fmt.Printf("failed to open issue %s in browser", issue.GetHTMLURL())
			}
		}

		if li.LinkedPrs {
			for url, pr := range prs {
				fmt.Println(pr)
				if li.BrowsePrs {
					if err := browser.OpenURL(url); err != nil {
						fmt.Printf("failed to open PR %s in browser", url)
					}
				}
			}
		}

		if li.Batch > 0 {
			if counter == li.Batch {
				continueListing := prompter.YN("Do you want to continue listing issues?", true)
				if !continueListing {
					return false, nil
				}
				counter = 1
				continue
			}
			counter++
		}
	}

	return true, nil
}

func crossReferencedPRs(repo gh.Repo, issueNumber int, prState string) (map[string]string, error) {
	prs := make(map[string]string, 0)

	timelines, err := repo.ListIssueTimeline(issueNumber)
	if err != nil {
		return nil, err
	}

	for _, timeline := range timelines {
		if timeline.GetEvent() == "cross-referenced" {
			if source := timeline.GetSource(); source != nil {
				if issue := source.GetIssue(); issue != nil {
					if issue.IsPullRequest() {
						if prState != "" {
							if prState == "merged" {
								merged, err := repo.PullRequestIsMerged(issue.GetNumber())
								if err != nil {
									return nil, err
								}
								if !merged {
									continue
								}
							} else if prState == "open" {
								if issue.GetState() != "open" {
									continue
								}
							} else if prState == "closed" {
								if issue.GetState() != "closed" {
									continue
								}
							}
						}

						if issue.GetHTMLURL() != "" {
							pr := fmt.Sprintf("\t- PR: %d\t%s\t%s", issue.GetNumber(), issue.GetHTMLURL(), issue.GetTitle())
							prs[issue.GetHTMLURL()] = pr
						}
					}
				}
			}
		}
	}
	return prs, nil
}
