package openai

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// Client wraps the OpenAI client
type Client struct {
	client *openai.Client
}

// NewClient creates a new OpenAI client
func NewClient(baseURL, apiKey string) *Client {
	// Create config with custom base URL
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL

	return &Client{
		client: openai.NewClientWithConfig(config),
	}
}

// GenerateCommitMessage generates a commit message using the OpenAI API
func (c *Client) GenerateCommitMessage(promptTemplate, diff, model string) (string, error) {
	// Replace {diff} placeholder in prompt
	prompt := strings.Replace(promptTemplate, "{diff}", diff, 1)

	// Create chat completion request
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant that generates Git commit messages following the format specified in the user's prompt.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   500,
			Temperature: 0.3,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	message := strings.TrimSpace(resp.Choices[0].Message.Content)

	// Clean up the message - remove surrounding quotes if present
	message = strings.Trim(message, "\"")
	message = strings.Trim(message, "'")

	// Clean up any markdown code blocks
	message = strings.TrimPrefix(message, "```")
	message = strings.TrimSuffix(message, "```")
	message = strings.TrimSpace(message)

	return message, nil
}
