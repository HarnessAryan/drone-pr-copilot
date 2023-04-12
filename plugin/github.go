package plugin

import (
	"context"
	"fmt"

	"github.com/google/go-github/v41/github"
)

func postReviewComment(ctx context.Context, client *github.Client, owner, repo string, prNumber int, feedbackList []*Feedback) error {
	// Check if the PR exists (again)
	fmt.Println("owner: ", owner)
	fmt.Println("repo: ", repo)
	fmt.Println("prNumber: ", prNumber)
	pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNumber)
	if pr == nil {
		return fmt.Errorf("PR not found")
	}

	// Prepare the draft review comments
	var draftComments []*github.DraftReviewComment
	for _, feedback := range feedbackList {
		if feedback.LineNumber < 0 {
			continue
		}
		fmt.Println("filename: ", feedback.Filename)
		fmt.Println("line number: ", feedback.RelativeLineNumber)
		fmt.Println("message: ", feedback.Suggestion)
		comment := &github.DraftReviewComment{
			Path:     github.String(feedback.Filename),
			Position: github.Int(feedback.RelativeLineNumber),
			Body:     github.String(feedback.Suggestion),
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
