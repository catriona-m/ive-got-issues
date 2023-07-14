package gh

import (
	"context"
	"fmt"
)

func (r Repo) PullRequestIsMerged(prNumber int) (bool, error) {

	client := r.NewClient()
	isMerged, _, err := client.PullRequests.IsMerged(context.Background(), r.Owner, r.Name, prNumber)
	if err != nil {
		return false, fmt.Errorf("error checking if pull request %d is merged: %v", prNumber, err)
	}

	return isMerged, nil
}
