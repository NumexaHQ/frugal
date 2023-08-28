package utils

import (
	"encoding/json"
	"errors"
	"io"
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
func ExtractContentFromRequestBody(rb io.ReadCloser) (string, error) {
	var reqBody RequestBody

	// Read the request body
	requestBodyBytes, err := io.ReadAll(rb)
	if err != nil {
		return "", err
	}

	// Unmarshal the JSON request body into the RequestBody struct
	if err := json.Unmarshal(requestBodyBytes, &reqBody); err != nil {
		return "", err
	}

	// Check if there is at least one message in the messages array
	if len(reqBody.Messages) < 1 {
		return "", errors.New("no messages found in the request body")
	}

	// Extract the content from the first message
	content := reqBody.Messages[0].Content

	return content, nil
}
