package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v52/github"
)

func (r Repo) ListIssuesByRepo(opts github.IssueListByRepoOptions) ([]*github.Issue, int, error) {

	client := r.NewClient()
	issues, resp, err := client.Issues.ListByRepo(context.Background(), r.Owner, r.Name, &opts)

	// TODO
	if resp.StatusCode != 200 {

	}

	if err != nil {
		return []*github.Issue{}, 0, err
	}
	return issues, resp.NextPage, nil
}

func (r Repo) AddLabelsToIssue(labels []string, issueNumber int) error {
	client := r.NewClient()

	_, _, err := client.Issues.AddLabelsToIssue(context.Background(), r.Owner, r.Name, issueNumber, labels)

	if err != nil {
		return fmt.Errorf("adding labels %v to issue %d", labels, issueNumber)
	}
	return nil
}

func (r Repo) AddCommentToIssue(comment string, issueNumber int) error {
	client := r.NewClient()

	_, _, err := client.Issues.CreateComment(context.Background(), r.Owner, r.Name, issueNumber, &github.IssueComment{Body: &comment})

	if err != nil {
		return fmt.Errorf("adding comment to issue %d", issueNumber)
	}
	return nil
}

func (r Repo) ListIssueTimeline(issueNumber int) ([]*github.Timeline, error) {
	client := r.NewClient()
	timelines := make([]*github.Timeline, 0)

	opts := github.ListOptions{
		Page:    1,
		PerPage: 100,
	}
	for {
		timeline, resp, err := client.Issues.ListIssueTimeline(context.Background(), r.Owner, r.Name, issueNumber, &opts)
		if err != nil {
			return nil, fmt.Errorf("requesting timeline for issue %d : %v", issueNumber, err)
		}
		timelines = append(timelines, timeline...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage

	}
	return timelines, nil
}
