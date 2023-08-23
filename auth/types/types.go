package types

import (
	openai "github.com/sashabaranov/go-openai"
)

// Response struct
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   openai.Usage           `json:"usage"`
}

type ChatCompletionChoice struct {
	Message      openai.ChatCompletionMessage `json:"message"`
	FinishReason string                       `json:"finish_reason"`
}

type MetricsResponse struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
