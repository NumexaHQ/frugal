package utils

import (
	"encoding/json"
	"errors"
)

var SensitiveHeaders = map[string]bool{
	"X-Numexa-Api-Key": true,
	"Authorization":    true,
	"Organization":     true,
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Helper function to extract content from the request body
func ExtractContentFromRequestBody(requestBodyString string) (string, error) {
	var reqBody RequestBody

	// Unmarshal the JSON request body into the RequestBody struct
	if err := json.Unmarshal([]byte(requestBodyString), &reqBody); err != nil {
		return "", err
	}

	// Check if there is at least one message in the messages array
	if len(reqBody.Messages) < 1 {
		return "", errors.New("No messages found in the request body")
	}

	// Extract the content from the first message
	content := reqBody.Messages[0].Content

	return content, nil
}
