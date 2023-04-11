package plugin

import (
	"context"
	"encoding/json"
	"fmt"

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

func (c *client) Feedback(fileDiffs []*FileDiff) []*Feedback {
	feedback := []*Feedback{}
	prompt := `
	Your job is to review pull request changes in code and return back improvements based on best practices that you can find online.
	I will give you as input the file before and after the changes have been made to it.
	Your job is to review what can be improved and return back. Your reply needs to be just a valid json and nothing else.
	The response should be a list of objects which contain 'line_number' and 'suggestion'.
	The suggestion should be concise and to the point. Here is the original file:
	Here is the original file: %s
	and here is the new file: %s
	`
	for _, diff := range fileDiffs {
		var old string
		var new string
		for _, line := range diff.PreviousLines {
			old += fmt.Sprintf("%d %s\n", line.Number, line.Content)
		}
		for _, line := range diff.NewLines {
			new += fmt.Sprintf("%d %s\n", line.Number, line.Content)
		}
		resp, err := c.client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: fmt.Sprintf("%s %s %s", prompt, old, new),
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
		}

		feedback = append(feedback, f...)

		fmt.Println("successfully received feedback from openai")
	}

	return feedback
}
