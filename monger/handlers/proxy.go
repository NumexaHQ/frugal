package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/NumexaHQ/monger/model"
	nxopenaiModel "github.com/NumexaHQ/monger/model/openai"
	"github.com/NumexaHQ/monger/utils"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/sirupsen/logrus"
)

func (h *Handler) OpenAIProxy(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now()
	originalURL := r.URL.String()

	index := strings.Index(originalURL, "/v1/openai/")
	if index == -1 {
		http.Error(w, "Invalid URL", http.StatusNotFound)
		return
	}

	endpoint := originalURL[index+len("/v1/openai/"):]
	newURL := "https://api.openai.com/v1/" + endpoint
	logrus.Infof("New URL: %s", newURL)
	useCache := r.Header.Get("X-Numexa-Cache") == "true"
	// new url
	u, err := url.Parse(newURL)
	if err != nil {
		http.Error(w, "Error parsing URL", http.StatusInternalServerError)
		return
	}

	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	r.Body = rdr2

	// Extract prompt from request body
	requestBodyString := string(buf)
	prompt, err := utils.ExtractContentFromRequestBody(requestBodyString)
	if err != nil {
		logrus.Errorf("Error extracting content from request body: %v", err)
		return
	}

	if useCache {
		cacheResponse, cacheErr := h.fetchFromCache(prompt)
		if cacheErr != nil {
			h.serveCachedResponse(w, r, cacheResponse)
			return
		}
	}

	// Create a new proxy request with the target URL
	proxyReq, err := http.NewRequest(r.Method, u.String(), rdr1)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating proxy request: %s", err), http.StatusInternalServerError)
		return
	}

	// Copy all headers from the original request to the proxy request
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	apiKey := r.Header.Get("X-Numexa-Api-Key")

	pr, err := model.ProxyRequestBuilderForHTTPRequest(r, requestTime, h.AuthDB, newURL, apiKey)
	if err != nil {
		logrus.Errorf("Error building proxy request: %v", err)
	}

	go func() {
		h.ChConfig.ReqC <- &pr
	}()

	// Create a new CustomResponseWriter to capture the response data
	customWriter := &model.CustomResponseWriter{ResponseWriter: w}

	// Use the default transport (http.DefaultTransport) for the proxy request
	proxyClient := http.DefaultClient

	// Create a new retryable http client
	retry := r.Header.Get("X-Numexa-Retry")
	retryCount := r.Header.Get("X-Numexa-Retry-Count")
	if retry == "true" && retryCount != "" {
		retryClient := retryablehttp.NewClient()
		retryCountInt, err := strconv.Atoi(retryCount)
		if err != nil {
			logrus.Errorf("Error converting retry count to int: %v", err)
			retryCountInt = 3
		}

		retryClient.RetryMax = retryCountInt
		proxyClient = retryClient.StandardClient() // *http.Client
	}

	initTime := time.Now()

	// Perform the proxy request and get the response
	proxyResp, err := proxyClient.Do(proxyReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error performing proxy request: %s", err), http.StatusInternalServerError)
		return
	}
	defer proxyResp.Body.Close()

	// Copy all headers from the proxy response to the original response
	for key, values := range proxyResp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code for the original response to match the proxy response
	w.WriteHeader(proxyResp.StatusCode)

	// Copy the response body from the proxy response to the original response
	io.Copy(customWriter, proxyResp.Body)

	responseTime := time.Now()

	responseBody := customWriter.Body.Bytes()

	proxyResp.Body = io.NopCloser(bytes.NewBuffer(responseBody))
	proxyResp.ContentLength = int64(len(responseBody))
	var reqBody nxopenaiModel.RequestBody
	err = json.Unmarshal([]byte(pr.RequestBody), &reqBody)
	if err != nil {
		logrus.Errorf("Error unmarshalling request body: %v", err)
		return
	}
	tokenLength := 0
	if reqBody.Stream {
		tokenLength, err = utils.GetBPETokenSizeByModel(reqBody.Messages[0].Content, reqBody.Model)
		if err != nil {
			logrus.Errorf("Error getting token length: %v", err)
			return
		}
	}
	pRes, err := model.ProxyResponseBuilderForHTTPResponse(r.Context(), proxyResp, h.AuthDB, initTime, tokenLength, responseTime, apiKey)
	if err != nil {
		logrus.Errorf("Error building proxy response: %v", err)
		return
	}
	if useCache {
		prompt := prompt
		h.storeInCache(prompt, r, proxyResp, pRes.ResponseBody)
	}
	h.ChConfig.ResC <- &pRes
}

// caching functions
func (h *Handler) storeInCache(prompt string, r *http.Request, proxyResp *http.Response, responseBody string) {
	// Create a map to represent the data to be stored in the cache
	cacheData := map[string]interface{}{
		"prompt": prompt,
		"answer": responseBody, // You can modify this part as needed
	}

	// Marshal the cache data to JSON
	cacheJSON, err := json.Marshal(cacheData)
	if err != nil {
		logrus.Errorf("Error marshaling cache data to JSON: %v", err)
		return
	}

	// Create a new request to store data in the cache server
	cacheURL := "http://nxa-cache:8000/put"
	cacheReq, err := http.NewRequest(http.MethodPost, cacheURL, bytes.NewBuffer(cacheJSON))

	if err != nil {
		logrus.Errorf("Error creating cache request: %v", err)
		return
	}

	// Set headers for the cache request
	cacheReq.Header.Set("Accept", "application/json")
	cacheReq.Header.Set("Content-Type", "application/json")
	// Add any additional headers you need here

	// Perform the cache request to store the response
	_, cacheErr := http.DefaultClient.Do(cacheReq)
	if cacheErr != nil {
		logrus.Errorf("Error storing response in cache: %v", cacheErr)
	}

	// Reset the response body for subsequent use
	proxyResp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte{}))
}

func (h *Handler) serveCachedResponse(w http.ResponseWriter, r *http.Request, cacheResp *http.Response) {
	// Copy headers from the cached response to the original response
	for key, values := range cacheResp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code for the original response to match the cached response
	w.WriteHeader(cacheResp.StatusCode)

	// Copy the response body from the cached response to the original response
	io.Copy(w, cacheResp.Body)

	// Close the cached response body
	cacheResp.Body.Close()
}

func (h *Handler) fetchFromCache(prompt string) (*http.Response, error) {
	// Create a map to represent the request data
	requestData := map[string]interface{}{
		"prompt": prompt,
	}

	// Marshal the request data to JSON
	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	// Create a new request to fetch data from the cache server
	cacheURL := "http://nxa-cache:8000/get"
	cacheReq, err := http.NewRequest(http.MethodPost, cacheURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}

	// Set headers for the cache request
	cacheReq.Header.Set("Content-Type", "application/json")
	// Add any additional headers you need here

	// Perform the cache request
	cacheResp, err := http.DefaultClient.Do(cacheReq)
	if err != nil {
		return nil, err
	}

	return cacheResp, nil
}
