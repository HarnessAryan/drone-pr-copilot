package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type client struct {
	token  string
	client *openai.Client
}

func New(opts ...Option) *client {
	c := new(client)
	for _, option := range opts {
		option(c)
	}
	c.client = openai.NewClient(c.token)
	return c
}

func (c *client) Feedback(ctx context.Context, fileDiffs []*FileDiff) []*Feedback {
	feedback := []*Feedback{}
	prompt := `
Your job is to review pull request changes in Golang code and return back improvements based on best practices found in the official Golang documentation, style guides, or other reputable online sources.
I will give you as input the file before and after the changes have been made to it.
Your job is to review what can be improved and return back. Your reply needs to be just a valid JSON and nothing else.
The response should be a list of objects which contain 'line_number' and 'suggestion'.
The suggestion should be concise and to the point. Please follow the official Golang guidelines and best practices, such as those in "Effective Go" and the Go Code Review Comments, while providing your suggestions.
Here is the original file:
%s
And here is the new file:
%s`

	for _, diff := range fileDiffs {
		var old string
		var new string
		for _, line := range diff.PreviousLines {
			old += fmt.Sprintf("%d %s\n", line.Number, line.Content)
		}
		for _, line := range diff.NewLines {
			new += fmt.Sprintf("%d %s\n", line.Number, line.Content)
		}
		cont := fmt.Sprintf(prompt, old, new)
		fmt.Println("content is: ", cont)
		resp, err := c.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: cont,
					},
				},
			},
		)
		if err != nil {
			fmt.Println("ChatCompletion error: %w", err)
			continue
		}
		// Try to unmarshal the response into Feedback struct
		content := resp.Choices[0].Message.Content
		var f []*Feedback
		err = json.Unmarshal([]byte(content), &f)
		if err != nil {
			fmt.Printf("could not unmarshal response: %s\n", content)
			continue
		}

		for _, entry := range f {
			entry.Filename = diff.Name
			entry.RelativeLineNumber = findInDiff(diff.NewLines[entry.LineNumber-1].Content, diff.Diff)
		}

		feedback = append(feedback, f...)

	}

	fmt.Printf("received %d feedback!", len(feedback))

	return feedback
}

func findInDiff(s string, diff []Line) int {
	for _, k := range diff {
		if sanitize(s) == sanitize(k.Content) {
			return k.Number - 1
		}
	}
	return -1
}

func sanitize(s string) string {
	if strings.HasPrefix(s, "+") {
		s = s[1:]
	}
	if strings.HasPrefix(s, "-") {
		s = s[1:]
	}
	s = strings.TrimSpace(s)
	return s
}
