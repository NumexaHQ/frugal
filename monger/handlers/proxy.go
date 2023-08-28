package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	commonConstants "github.com/NumexaHQ/captainCache/numexa-common/constants"
	"github.com/NumexaHQ/monger/model"
	nxopenaiModel "github.com/NumexaHQ/monger/model/openai"
	gptcache "github.com/NumexaHQ/monger/pkg/cache"
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

	endpoint := originalURL[index+len(commonConstants.OPENAI_PROXY_STUB):]
	newURL := commonConstants.OPENAI_BASE_URL + endpoint
	useCache := r.Header.Get("X-Numexa-Cache") == "true"

	u, err := url.Parse(newURL)
	if err != nil {
		http.Error(w, "Error parsing URL", http.StatusInternalServerError)
		return
	}

	buf, _ := io.ReadAll(r.Body)
	rdr1 := io.NopCloser(bytes.NewBuffer(buf))
	rdr2 := io.NopCloser(bytes.NewBuffer(buf))
	r.Body = rdr2

	// Extract prompt from request body
	prompt, err := utils.ExtractContentFromRequestBody(io.NopCloser(bytes.NewBuffer(buf)))
	if err != nil {
		logrus.Errorf("Error extracting content from request body: %v", err)
		return
	}

	gptCache := gptcache.New(prompt, useCache)

	// check if cache is enabled and if there is a cached response
	if gptCache.Enabled && gptCache.GetCachedAnswer() != "" {
		h.serveCachedAnswer(w, r, gptCache.GPTResponse)
		return
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

	// cache only if the response is 200
	if gptCache.Enabled && pRes.ResponseStatusCode == 200 {
		// if code is here, it means that the cache is enabled and there is no cached response yet
		err = gptCache.SetCachedAnswer(pRes.ResponseBody)
		if err != nil {
			logrus.Errorf("Error setting cached answer: %v", err)
			// donot return here, we still want to ingest the response
		}
	}

	h.ChConfig.ResC <- &pRes
}

func (h *Handler) serveCachedAnswer(w http.ResponseWriter, r *http.Request, gc gptcache.GPTCache) {
	// set content type
	w.Header().Set("Content-Type", "application/json")

	// set status code
	w.WriteHeader(200)

	w.Write([]byte(gc.Answer))
}
