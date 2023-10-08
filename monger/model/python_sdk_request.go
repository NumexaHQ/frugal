package model

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

type PythonSdkRequest struct {
	RequestTime    string                 `json:"request_time"`
	SourceIp       string                 `json:"source_ip"`
	RequestMethod  string                 `json:"request_method"`
	RequestUrl     string                 `json:"request_url"`
	RequestHeaders map[string]string      `json:"request_headers"`
	RequestBody    map[string]interface{} `json:"request_body"`
}

// CreateNewRequest Helper function to Create new Request received from python Sdk
func CreateNewRequest(err error, body []byte, r *http.Request) (error, *http.Request, time.Time) {
	bodyFromPythonSdk := &PythonSdkRequest{}
	err = json.Unmarshal(body, bodyFromPythonSdk)
	if err != nil {
		return err, r, time.Now()
	}
	// Creating Body Bytes for New Request
	newRequestBodyBytes, err := json.Marshal(bodyFromPythonSdk.RequestBody)
	if err != nil {
		return err, r, time.Now()
	}
	// Creating Body From Body Bytes For New Request
	newRequestBody := io.NopCloser(bytes.NewBuffer(newRequestBodyBytes))

	u, err := url.Parse(bodyFromPythonSdk.RequestUrl)
	if err != nil {
		return err, r, time.Now()
	}
	requestTimeStamp, err := time.Parse("2006-01-02 15:04:05.000", bodyFromPythonSdk.RequestTime)
	if err != nil {
		return err, r, time.Now()
	}
	// Setting up New Request With Existing Request Context
	newReq, err := http.NewRequestWithContext(r.Context(), r.Method, u.String(), newRequestBody)
	if err != nil {
		return err, newReq, time.Now()
	}
	//Setting Up new Request URl and RemoteAddress
	newReq.URL = u
	newReq.RemoteAddr = bodyFromPythonSdk.SourceIp

	// Copying Existing Request Headers into New Request
	for key, values := range r.Header {
		for _, value := range values {
			newReq.Header.Add(key, value)
		}
	}
	return nil, newReq, requestTimeStamp
}
