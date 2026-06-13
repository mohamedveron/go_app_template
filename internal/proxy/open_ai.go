package proxy

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	token string
}

func (ai *OpenAI) GetMessage(token string) string {
	client := openai.NewClient(ai.token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: token,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	str := fmt.Sprintf("%s", resp.Choices[0].Message.Content)
	return str
}

func NewOpenAI(token string) OpenAI {
	return OpenAI{token: token}
}
