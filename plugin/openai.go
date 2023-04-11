package plugin

import (
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

func (c *client) Feedback(files []*File) []*Feedback {
	return []*Feedback{}
}
