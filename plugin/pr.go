package plugin

import (
	"context"
	"fmt"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

func main() {
	token := "your_token_here"
	owner := "your_owner_here"
	repo := "your_repo_here"
	pullRequestNumber := 0 // Replace with the pull request number

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	pr, _, err := client.PullRequests.Get(ctx, owner, repo, pullRequestNumber)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	printPullRequestDetails(pr)
}

func printPullRequestDetails(pr *github.PullRequest) {
	fmt.Printf("Title: %s\n", pr.GetTitle())
	fmt.Printf("Author: %s\n", pr.GetUser().GetLogin())
	fmt.Printf("State: %s\n", pr.GetState())
	fmt.Printf("Created At: %s\n", pr.GetCreatedAt().String())
	fmt.Printf("URL: %s\n", pr.GetHTMLURL())
}
