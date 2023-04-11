package plugin

import "context"
import "github.com/google/go-github/v41/github"

func postReviewComment(ctx context.Context, client *github.Client, owner, repo string, prNumber int, feedback Feedback) error {
	comment := &github.DraftReviewComment{
		Path:     github.String(feedback.Filename),
		Position: github.Int(feedback.LineNumber),
		Body:     github.String(feedback.Message),
	}

	pullRequestReviewRequest := &github.PullRequestReviewRequest{
		Comments: []*github.DraftReviewComment{comment},
	}

	_, _, err := client.PullRequests.CreateReview(ctx, owner, repo, prNumber, pullRequestReviewRequest)
	return err
}
