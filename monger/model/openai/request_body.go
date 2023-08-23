package model

type RequestBody struct {
	N           int       `json:"n"`
	Model       string    `json:"model"`
	Stream      bool      `json:"stream"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature int       `json:"temperature"`
}
