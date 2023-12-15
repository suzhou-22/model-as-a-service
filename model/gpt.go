package model

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type Gpt struct {
	client *openai.Client
	model  string
}

func NewGptModel(token, baseUrl, model string) *Gpt {
	if model == "" {
		model = openai.GPT3Dot5Turbo
	}

	clientConfig := openai.DefaultConfig(token)
	clientConfig.BaseURL = baseUrl
	return &Gpt{openai.NewClientWithConfig(clientConfig), model}
}

func (c Gpt) Complete(context context.Context, prompt string) string {
	resp, err := c.client.CreateChatCompletion(
		context,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}

	return resp.Choices[0].Message.Content
}

var _ Interface = Gpt{}
