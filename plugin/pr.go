package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v41/github"
)

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

func GetFileDiff(ctx context.Context, client *github.Client, owner string, repo string, pullRequestNumber int) ([]FileDiff, error) {
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

	fileDiffs := []FileDiff{}

	for _, file := range files {
		name := file.GetFilename()
		var previousLines, newLines []Line

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

		fileDiffs = append(fileDiffs, FileDiff{
			Name:          name,
			PreviousLines: previousLines,
			NewLines:      newLines,
		})
	}

	// TODO: this can be deleted - its just for testing
	jsonOutput, _ := json.MarshalIndent(fileDiffs, "", "  ")
	fmt.Println(string(jsonOutput))

	return fileDiffs, nil
}
