// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

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
	tokens := strings.Split(args.Pipeline.Commit.Link, "/")
	prs := tokens[len(tokens)-1]
	pr, err := strconv.Atoi(prs)
	if err != nil {
		log.Fatalf("could not parse pr number")
	}
	fmt.Println("pr number: ", pr)
	fmt.Println("pipeline repo name: ", args.Pipeline.Repo.Name)
	githubClient := createGithubClient(ctx, args)

	fileDiffs, err := GetFileDiff(ctx, githubClient, args.Pipeline.Repo.Namespace, args.Pipeline.Repo.Name, pr)
	if err != nil {
		log.Fatalf("could not get file diff, err: %s", err)
	}

	// Pass in the file diffs to get feedback on them
	openAIClient := New(WithToken(args.OpenAIKey))
	feedback := openAIClient.Feedback(ctx, fileDiffs)

	err = postReviewComment(ctx, githubClient, args.Pipeline.Repo.Namespace, args.Pipeline.Repo.Name, pr, feedback)
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
