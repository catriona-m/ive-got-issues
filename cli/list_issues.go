package cli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Songmu/prompter"
	"github.com/google/go-github/v52/github"
	c "github.com/gookit/color"
	"github.com/ivegotissues/lib/gh"
	"github.com/pkg/browser"
)

type ListIssues struct {
	Labels     []string
	Content    string
	IssueState string
	LinkedPrs  bool
	PrState    string
	OpenIssues bool
	OpenPrs    bool
	Batch      int
	Token      string
	Repo       string
}

func (li ListIssues) ListIssues() error {

	opts := github.IssueListByRepoOptions{
		State: li.IssueState,
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	if len(li.Labels) > 0 {
		opts.Labels = li.Labels
	}

	repo := gh.NewRepo(li.Repo, li.Token)
	issueCount := 0
	totalIssuesCount := 0
	batchCounter := 1
	var continueListing bool
	for {

		issues, nextPage, err := repo.ListIssuesByRepo(opts)
		if err != nil {
			return fmt.Errorf("retrieving issues from github from page %d: %v", opts.ListOptions.Page, err)
		}
		continueListing, issueCount, err = li.listFilteredIssues(repo, issues, batchCounter)
		if err != nil {
			return err
		}

		totalIssuesCount += issueCount

		if li.Batch > 0 && totalIssuesCount+1 >= li.Batch {
			batchCounter = (totalIssuesCount + 1) % li.Batch
			if batchCounter == 0 {
				batchCounter = li.Batch
			}
		} else {
			batchCounter = totalIssuesCount + 1
		}

		if !continueListing || nextPage == 0 {
			c.Info.Printf("Finished listing %d issues.\n", totalIssuesCount)
			break
		}
		opts.ListOptions.Page = nextPage
	}

	return nil
}

func (li ListIssues) listFilteredIssues(repo gh.Repo, issues []*github.Issue, batchCounter int) (bool, int, error) {

	var err error
	issuesCount := 0

	re := &regexp.Regexp{}
	if li.Content != "" {
		re, err = regexp.Compile(li.Content)
		if err != nil {
			return false, 0, err
		}
	}

	for _, issue := range issues {
		if issue.IsPullRequest() {
			continue
		}

		if li.Content != "" {
			if !re.MatchString(issue.GetBody()) {
				continue
			}
		}

		prs := make(map[string]string)
		if li.LinkedPrs {
			prs, err = crossReferencedPRs(repo, issue.GetNumber(), li.PrState)
			if err != nil {
				return false, issuesCount, err
			}
			if len(prs) == 0 {
				continue
			}
		}

		date := strings.Split(issue.GetCreatedAt().String(), " ")[0]
		c.Printf("\n<cyan>#%d</> @ <green>%s</>\t%s\t%s\n", issue.GetNumber(), date, issue.GetTitle(), issue.GetHTMLURL())
		issuesCount++

		if li.OpenIssues {
			if err := browser.OpenURL(issue.GetHTMLURL()); err != nil {
				fmt.Printf("failed to open issue %s in browser", issue.GetHTMLURL())
			}
		}

		if li.LinkedPrs {
			for url, pr := range prs {
				c.Println(pr)
				if li.OpenPrs {
					if err := browser.OpenURL(url); err != nil {
						fmt.Printf("failed to open PR %s in browser", url)
					}
				}
			}
		}

		if li.Batch > 0 {
			if batchCounter == li.Batch {
				continueListing := prompter.YN("Do you want to continue listing issues?", true)
				if !continueListing {
					return false, issuesCount, nil
				}
				batchCounter = 1
				continue
			}
			batchCounter++
		}
	}

	return true, issuesCount, nil
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
						colour := ""

						switch issue.GetState() {
						case "open":
							if prState == "open" || prState == "" {
								colour = "<lightGreen>"
							} else {
								continue
							}

						case "closed":
							if prState != "" && prState != "merged" && prState != "closed" {
								continue
							}
							merged, err := repo.PullRequestIsMerged(issue.GetNumber())
							if err != nil {
								return nil, err
							}
							if merged {
								if prState == "" || prState == "merged" {
									colour = "<lightMagenta>"
								} else {
									continue
								}
							} else if prState == "" || prState == "closed" {
								colour = "<lightRed>"
							} else {
								continue
							}
						}

						if issue.GetHTMLURL() != "" {
							pr := c.Sprintf("\t%s#%d\t%s\t%s</>", colour, issue.GetNumber(), issue.GetTitle(), issue.GetHTMLURL())
							prs[issue.GetHTMLURL()] = pr
						}
					}
				}
			}
		}
	}
	return prs, nil
}
