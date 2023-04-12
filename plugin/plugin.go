// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"encoding/json"
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

func getFileContentAtCommit(ctx context.Context, client *github.Client, owner, repo, path, commitSHA string) (string, error) {
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{Ref: commitSHA})
	if err != nil {
		return "", err
	}

	decoded, err := fileContent.GetContent()
	if err != nil {
		return "", err
	}

	return decoded, nil
}

func convertContentToLines(content string) []Line {
	lines := strings.Split(content, "\n")
	lineStructs := make([]Line, len(lines))

	for i, line := range lines {
		lineStructs[i] = Line{
			Number:  i + 1,
			Content: line,
		}
	}

	return lineStructs
}

func GetFileDiff(ctx context.Context, client *github.Client, owner string, repo string, pullRequestNumber int) ([]*FileDiff, error) {
	pr, _, err := client.PullRequests.Get(ctx, owner, repo, pullRequestNumber)
	if err != nil {
		return nil, err
	}

	baseCommitID := pr.GetBase().GetSHA()
	latestCommitID := pr.GetHead().GetSHA()

	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, pullRequestNumber, nil)
	if err != nil {
		return nil, err
	}

	commits, _, err := client.PullRequests.ListCommits(ctx, owner, repo, pullRequestNumber, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	commitNumbers := make(map[string]int)
	for i, commit := range commits {
		commitNumbers[commit.GetSHA()] = i + 1
	}

	fileDiffs := []*FileDiff{}

	for _, file := range files {
		name := file.GetFilename()
		var previousLines, newLines []Line
		diff := file.GetPatch()
		commitNumber := commitNumbers[latestCommitID]

		beforePR, err := getFileContentAtCommit(ctx, client, owner, repo, name, baseCommitID)
		if err != nil {
			if !strings.Contains(err.Error(), "404 Not Found") {
				fmt.Printf("Error getting file content before PR: %v\n", err)
				continue
			}
		} else {
			previousLines = convertContentToLines(beforePR)
		}

		afterPR, err := getFileContentAtCommit(ctx, client, owner, repo, name, latestCommitID)
		if err != nil {
			fmt.Printf("Error getting file content after PR: %v\n", err)
			continue
		}
		newLines = convertContentToLines(afterPR)

		fileDiffs = append(fileDiffs, &FileDiff{
			Name:          name,
			PreviousLines: previousLines,
			NewLines:      newLines,
			Diff:          convertContentToLines(diff),
			CommitNumber:  commitNumber,
		})
	}

	_, err = json.MarshalIndent(fileDiffs, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	return fileDiffs, nil
}

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
