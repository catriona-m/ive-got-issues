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
