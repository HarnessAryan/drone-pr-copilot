// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
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
	// Printing some data that we will need
	fmt.Println("pipeline namespace: ", args.Pipeline.Repo.Namespace)
	fmt.Println("pipeline name: ", args.Pipeline.Repo.Name)
	fmt.Println("pr number: ", args.Pipeline.PullRequest.Number)
	fmt.Println("commit author name: ", args.Pipeline.Commit.Author.Name)
	fmt.Println("pipeline repo name: ", args.Pipeline.Repo.Name)
	githubClient := createGithubClient(ctx, args)
	// feedbackList := []*Feedback{
	// 	{
	// 		Filename:   "renovate.json",
	// 		LineNumber: 1,
	// 		Suggestion: "Replace 'fmt.Println()' with 'log.Println()'",
	// 		Message:    "Use 'log' package instead of 'fmt' for better control over logging output.",
	// 		Severity:   "warning",
	// 	},
	// 	{
	// 		Filename:   "renovate.json",
	// 		LineNumber: 4,
	// 		Suggestion: "Add error handling for the function call",
	// 		Message:    "Error handling is missing for the function call. It's important to handle errors to avoid unexpected behavior.",
	// 		Severity:   "error",
	// 	},
	// }

	fileDiffs, err := GetFileDiff(ctx, githubClient, args.Pipeline.Repo.Namespace, args.Pipeline.Repo.Name, args.Pipeline.PullRequest.Number)
	if err != nil {
		log.Fatalf("could not get file diff, err: %s", err)
	}

	// Pass in the file diffs to get feedback on them
	openAIClient := New(WithToken(args.OpenAIKey))
	feedback := openAIClient.Feedback(ctx, fileDiffs)

	err = postReviewComment(ctx, githubClient, args.Pipeline.Commit.Author.Name, args.Pipeline.Repo.Name, args.Pipeline.PullRequest.Number, feedback)
	if err != nil {
		log.Fatalf("could not post review comments, err: %s", err)
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
