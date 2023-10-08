package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andybalholm/brotli"
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

// EncodeIntoBrotli helper function to encode the marshal data into brotli format
func EncodeIntoBrotli(newResponseBodyBytes []byte) ([]byte, error) {
	brotliBytes := &bytes.Buffer{}
	writer := brotli.NewWriter(brotliBytes)
	_, err := writer.Write(newResponseBodyBytes)
	if err != nil {
		return nil, fmt.Errorf("error Encoding into brotli %s", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error while closing brotli Writer %s", err)
	}
	return brotliBytes.Bytes(), err
}

// EncodeIntoGzip helper function to encode the marshal data into gzip format
func EncodeIntoGzip(data []byte) ([]byte, error) {
	gzipBytes := &bytes.Buffer{}
	writer := gzip.NewWriter(gzipBytes)
	_, err := writer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error Encoding into Gzip %s", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error while closing gzip Writer %s", err)
	}
	return gzipBytes.Bytes(), nil
}
