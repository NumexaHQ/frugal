package model

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	nxAuthDB "github.com/NumexaHQ/captainCache/pkg/db"
	openaiModel "github.com/NumexaHQ/monger/model/openai"
	"github.com/NumexaHQ/monger/utils"
	"github.com/andybalholm/brotli"
	log "github.com/sirupsen/logrus"
)

// CustomResponseWriter is a custom ResponseWriter to capture the response data
type CustomResponseWriter struct {
	http.ResponseWriter
	Status int
	Header http.Header
	Body   bytes.Buffer
}

// WriteHeader overrides the WriteHeader method to capture the status code
func (c *CustomResponseWriter) WriteHeader(status int) {
	c.Status = status
	c.ResponseWriter.WriteHeader(status)
}

// Write overrides the Write method to capture the response body
func (c *CustomResponseWriter) Write(b []byte) (int, error) {
	c.ResponseWriter.Write(b)
	return c.Body.Write(b)
}

type ProxyResponse struct {
	RequestID          string `json:"request_id"`
	InitiatedTimestamp int64  `json:"initiated_timestamp"`
	ResponseTimestamp  int64  `json:"response_timestamp"`
	ResponseStatusCode uint16 `json:"response_status_code"`
	ResponseHeaders    string `json:"response_headers"`
	ResponseBody       string `json:"response_body"`
	Provider           string `json:"provider"`
	UserID             int32  `json:"user_id"`
	ProjectID          int32  `json:"project_id"`
}

func (p *ProxyResponse) SetUserIdentifier(ctx context.Context, authDB nxAuthDB.DB, apiKey string) error {
	// using the nxaAPIToken, get the user_id and project_id
	log.Infof("Setting user identifier for API key: %v", apiKey)
	numexaAPIKeyObj, err := authDB.GetAPIkeyByApiKey(ctx, apiKey)
	if err != nil {
		log.Errorf("Error getting API key: %v", err)
		return err
	}
	// set the user_id and project_id
	p.UserID = numexaAPIKeyObj.UserID
	p.ProjectID = numexaAPIKeyObj.ProjectID
	return nil
}

func ProxyResponseBuilderForHTTPResponse(ctx context.Context, r *http.Response, authDB nxAuthDB.DB, initTime time.Time, promptTokenLen int, responseTime time.Time, apiKey string) (ProxyResponse, error) {
	// header to map
	header := make(map[string]string)
	for k, v := range r.Header {
		// skip sensitive headers
		_, ok := utils.SensitiveHeaders[k]
		if ok {
			continue
		}
		header[k] = v[0]
	}

	body := ""
	var err error
	if r.Body != nil {
		contentEncoding := r.Header.Get("Content-Encoding")
		contentType := header["Content-Type"]
		log.Infof("Content encoding: %v", contentEncoding)
		if contentEncoding == "gzip" {
			log.Info("Content encoding is gzip, skipping")
			body, err = ungzip(r.Body)
			if err != nil {
				log.Errorf("Error ungziping: %v", err)
				return ProxyResponse{}, err
			}
			// todo
		} else if contentEncoding == "br" {
			log.Info("Content encoding is brotli, unbrotliing")
			body, err = unbrotli(r.Body)
			if err != nil {
				log.Errorf("Error unbrotliing: %v", err)
				return ProxyResponse{}, err
			}
		} else {
			if contentType == "text/event-stream" {
				log.Info("Content type is text/event-stream, unchunking")
				b, err := parseStream(r.Body)
				if err != nil {
					log.Errorf("Error parsing stream: %v", err)
					return ProxyResponse{}, err
				}
				b.Usage.PromptTokens = promptTokenLen
				b.Usage.TotalTokens = promptTokenLen + b.Usage.CompletionTokens
				bodyString, err := json.MarshalIndent(b, "", "  ")
				body = string(bodyString)
				if err != nil {
					log.Errorf("Error marshalling body: %v", err)
					return ProxyResponse{}, err
				}
			} else {
				buf := new(bytes.Buffer)
				buf.ReadFrom(r.Body)
				body = buf.String()
			}
		}
	}

	hB, err := json.Marshal(header)
	if err != nil {
		log.Errorf("Error marshalling header: %v", err)
		return ProxyResponse{}, err
	}

	// get requset idfrom context
	rid := GetRequestID(ctx)
	pr := ProxyResponse{
		RequestID:          rid,
		InitiatedTimestamp: initTime.Unix(),
		ResponseTimestamp:  responseTime.Unix(),
		ResponseStatusCode: uint16(r.StatusCode),
		ResponseHeaders:    string(hB),
		ResponseBody:       body,
		Provider:           "openapi",
	}

	err = pr.SetUserIdentifier(context.Background(), authDB, apiKey)
	if err != nil {
		log.Errorf("Error setting user identifier: %v", err)
		return ProxyResponse{}, err
	}
	return pr, nil
}

func unbrotli(b io.ReadCloser) (string, error) {
	reader := brotli.NewReader(b)
	// defer reader.Close()
	respBody, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("error decoding br response", err)
	}
	return string(respBody), nil
}

func ungzip(b io.ReadCloser) (string, error) {
	reader, err := gzip.NewReader(b)
	if err != nil {
		log.Fatal("error decoding gzip response", err)
	}
	defer reader.Close()
	respBody, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("error decoding gzip response", err)
	}
	return string(respBody), nil
}

func parseStream(b io.ReadCloser) (openaiModel.ResponseBody, error) {
	scanner := bufio.NewScanner(b)
	var responseBodyStream []openaiModel.ResponseBody

	for scanner.Scan() {
		chunck := scanner.Text()
		// this might be a hack, but it works for now
		// because marut is a peanut
		if chunck == "[DONE]" || chunck == "data: [DONE]" {
			break
		}
		// continue if empty line
		if chunck == "" {
			continue
		}
		// remove data: from the beginning of the line
		chunck = strings.Replace(chunck, "data: ", "", 1)
		var rb openaiModel.ResponseBody
		err := json.Unmarshal([]byte(chunck), &rb)
		if err != nil {
			log.Errorf("Error unmarshalling chunk: %v", err)
			return openaiModel.ResponseBody{}, err
		}
		responseBodyStream = append(responseBodyStream, rb)
	}

	nonStream := streamToNonStream(responseBodyStream)
	return nonStream, nil
}

func streamToNonStream(rbs []openaiModel.ResponseBody) openaiModel.ResponseBody {
	log.Printf("len %d", len(rbs))
	var rbFinal openaiModel.ResponseBody
	var msg string

	index := 1

	for _, rb := range rbs {
		rbFinal = openaiModel.ResponseBody{
			ID:      rb.ID,
			Object:  rb.Object,
			Created: rb.Created,
			Model:   rb.Model,
			Choices: []openaiModel.Choices{
				{
					Index:        rb.Choices[0].Index,
					Message:      openaiModel.Message{},
					FinishReason: rb.Choices[0].FinishReason,
				},
			},
		}
		msg += rb.Choices[0].Delta.Content
		index++
	}

	tokenLength, err := utils.GetBPETokenSizeByModel(msg, rbFinal.Model)
	if err != nil {
		log.Errorf("Error getting token length: %v", err)
	}
	// log.Infof("Token length: %d", tokenLength)

	rbFinal.Choices[0].Message.Content = msg
	rbFinal.Usage.CompletionTokens = tokenLength
	return rbFinal
}
