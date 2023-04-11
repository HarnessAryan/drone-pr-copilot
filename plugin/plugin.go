// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
	"log"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	// OpenAI key
	OpenAIKey   string `envconfig:"PLUGIN_OPENAI_KEY"`
	GithubToken string `envconfig:"PLUGIN_GITHUB_TOKEN"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	ctx = context.Background()
	githubClient := createGithubClient(ctx, args)
	feedbackList := []*Feedback{
		{
			Filename:   "renovate.json",
			LineNumber: 1,
			Suggestion: "Replace 'fmt.Println()' with 'log.Println()'",
			Message:    "Use 'log' package instead of 'fmt' for better control over logging output.",
			Severity:   "warning",
		},
		{
			Filename:   "renovate.json",
			LineNumber: 4,
			Suggestion: "Add error handling for the function call",
			Message:    "Error handling is missing for the function call. It's important to handle errors to avoid unexpected behavior.",
			Severity:   "error",
		},
	}

	err := postReviewComment(ctx, githubClient, args.Pipeline.Repo.Namespace, args.Pipeline.Repo.Name, args.Pipeline.PullRequest.Number, feedbackList)
	if err != nil {
		log.Fatalf("Error posting review comment: %v", err)
	}

	return nil
}

func createGithubClient(ctx context.Context, args Args) *github.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: args.GithubToken},
	)

	httpClient := oauth2.NewClient(ctx, tokenSource)

	// Create a new GitHub client
	return github.NewClient(httpClient)
}
