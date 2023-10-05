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

var OpenApiMandatoryHeaders = map[string]bool{
	"Authorization": true,
	"Content-Type":  true,
}

type Config struct {
	Mode    string   `json:"mode"`
	Options []Option `json:"options"`
}

type Option struct {
	Provider       string `json:"provider"`
	VirtualKey     string `json:"virtual_key"`
	OverrideParams Model  `json:"override_params"`
}

type Model struct {
	Model string `json:"model"`
}

type Params struct {
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type RequestBody struct {
	Config Config `json:"config"`
	Params Params `json:"params"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type FinalRequestBody struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

// Helper function to extract content from the request body
func ExtractContentFromRequestBody(rb io.ReadCloser) (string, error) {
	var reqBody FinalRequestBody

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
