package plugin

import (
	"context"
	"fmt"
)
import "github.com/google/go-github/v41/github"

func postReviewComment(ctx context.Context, client *github.Client, owner, repo string, prNumber int, feedbackList []*Feedback) error {
	// Check if the PR exists (again)
	pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNumber)
	if pr == nil {
		return fmt.Errorf("PR not found")
	}

	// Prepare the draft review comments
	var draftComments []*github.DraftReviewComment
	for _, feedback := range feedbackList {
		comment := &github.DraftReviewComment{
			Path:     github.String(feedback.Filename),
			Position: github.Int(feedback.LineNumber),
			Body:     github.String(feedback.Message),
		}
		draftComments = append(draftComments, comment)
	}

	// Prepare the pull request review request
	pullRequestReviewRequest := &github.PullRequestReviewRequest{
		Event:    github.String("REQUEST_CHANGES"),
		Body:     github.String("Please address the suggested inline changes."),
		Comments: draftComments,
	}

	// Create the review
	_, _, err = client.PullRequests.CreateReview(ctx, owner, repo, prNumber, pullRequestReviewRequest)
	return err
}
