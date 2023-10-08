package model

import (
	"bytes"
	"encoding/json"
	"github.com/NumexaHQ/monger/utils"
	"io"
	"net/http"
	"time"
)

type PythonSdkResponse struct {
	InitiatedTimestamp string                 `json:"initiated_timestamp"`
	ResponseTimestamp  string                 `json:"response_timestamp"`
	ResponseStatusCode uint16                 `json:"response_status_code"`
	ResponseBody       map[string]interface{} `json:"response_body"`
}

// CreateNewResponse Helper function to Create new Response received from python Sdk
func CreateNewResponse(r *http.Request, body []byte) (error, *http.Response, time.Time, time.Time) {
	pythonSdkResponse := &PythonSdkResponse{}
	err := json.Unmarshal(body, pythonSdkResponse)
	if err != nil {
		return err, &http.Response{}, time.Now(), time.Now()
	}
	// Creating Body Bytes for New Response
	newResponseBodyBytes, err := json.Marshal(pythonSdkResponse.ResponseBody)
	if err != nil {
		return err, &http.Response{}, time.Now(), time.Now()
	}
	contentEncoding := r.Header.Get("Content-Encoding")
	initiatedTimeStamp, err := time.Parse("2006-01-02 15:04:05.000", pythonSdkResponse.InitiatedTimestamp)
	if err != nil {
		return err, &http.Response{}, time.Now(), time.Now()
	}
	responseTimestamp, err := time.Parse("2006-01-02 15:04:05.000", pythonSdkResponse.ResponseTimestamp)
	if err != nil {
		return err, &http.Response{}, time.Now(), time.Now()
	}
	var bodyBytes []byte
	if contentEncoding == "br" {
		// Encoding Into Brotli
		brotliBytes, err := utils.EncodeIntoBrotli(newResponseBodyBytes)
		if err != nil {
			return err, &http.Response{}, time.Now(), time.Now()
		}
		bodyBytes = brotliBytes

	} else if contentEncoding == "gzip" {
		// Encoding Into Gzip
		gzipBytes, err := utils.EncodeIntoGzip(newResponseBodyBytes)
		if err != nil {
			return err, &http.Response{}, time.Now(), time.Now()
		}
		bodyBytes = gzipBytes
	} else {
		bodyBytes = newResponseBodyBytes
	}
	proxyResponse := &http.Response{
		StatusCode:    int(pythonSdkResponse.ResponseStatusCode),
		Body:          io.NopCloser(bytes.NewBuffer(bodyBytes)),
		ContentLength: int64(len(body)),
		Header:        r.Header,
	}
	return nil, proxyResponse, initiatedTimeStamp, responseTimestamp
}
