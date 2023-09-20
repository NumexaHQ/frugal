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

	authModel "github.com/NumexaHQ/captainCache/model"
	commonConstants "github.com/NumexaHQ/captainCache/numexa-common/constants"
	authConstants "github.com/NumexaHQ/captainCache/pkg/constants"
	"github.com/NumexaHQ/captainCache/pkg/providerkeys"
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
	userId := int32(0)
	gptCache := gptcache.Cache{}

	apiKey := r.Header.Get("X-Numexa-Api-Key")

	// check if openai key is present
	// todo: pprof profile this!!!
	if r.Header.Get("Organization") == "" || r.Header.Get("Authorization") == "" {
		// check if api key has associated openai key
		// if not, return error
		// else, set the openai key in the header
		//
		isValid, providerKey, keyProperty, providerSecrets, err := h.AuthDB.CheckProviderAndNXAKeyPropertyFromNXAKey(r.Context(), apiKey, authConstants.PROVIDER_OPENAI)
		if err != nil {
			logrus.Errorf("Error checking provider and nxa key property from nxa key: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !isValid {
			http.Error(w, "Invalid/Expired API key", http.StatusUnauthorized)
			return
		}

		logrus.Infof("provider key: %+v", providerKey)
		logrus.Infof("providerSecrets: %+v", providerSecrets)

		if providerKey.Provider != authConstants.PROVIDER_OPENAI {
			http.Error(w, "Invalid provider", http.StatusUnauthorized)
			return
		}

		// todo: use keyproperty, to enforce rules
		logrus.Debugf("key property: %v", keyProperty)

		keys := make(map[string]string)

		for _, secrets := range providerSecrets {
			keys[secrets.Type] = secrets.Key
		}

		// here kp.Keys is encrypted keys
		kp := authModel.ProviderKeys{
			Name:      providerKey.Name,
			Provider:  providerKey.Provider,
			Keys:      keys,
			ProjectId: providerKey.ProjectID,
		}

		kpB, err := json.Marshal(kp)
		if err != nil {
			logrus.Errorf("Error marshalling provider keys: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		provider, err := providerkeys.GetProvider(providerKey.Provider, kpB, true)
		if err != nil {
			logrus.Errorf("Error getting provider: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		decryptedKeys, err := provider.GetDecryptedKeys(r.Context(), h.AuthDB)
		if err != nil {
			logrus.Errorf("Error getting decrypted keys: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// set the openai key in the header
		r.Header.Set("Organization", decryptedKeys["openai_org"])
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", decryptedKeys["openai_key"]))
	}

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

	if useCache {
		numexaAPIKeyObj, err := h.AuthDB.GetAPIkeyByApiKey(r.Context(), apiKey)
		if err != nil {
			logrus.Errorf("Error getting API key: %v", err)
			return
		}
		userId = numexaAPIKeyObj.UserID

		gptCache = gptcache.New(prompt, userId, useCache)

		// check if cache is enabled and if there is a cached response
		if gptCache.GetCachedAnswer() != "" {
			h.serveCachedAnswer(w, r, gptCache.GPTResponse)
			gptCache.IngestCachedRequest(r, requestTime, h.AuthDB, newURL, apiKey, h.ChConfig)
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
	apiKey := r.Header.Get("X-Numexa-Api-Key")

	initTime := time.Now()

	w.Header().Set("Content-Type", "application/json")

	// set status code
	w.WriteHeader(200)

	w.Write([]byte(gc.Answer))

	responseTime := time.Now()

	//tokenLength
	tokenLength := 0

	pRes, err := model.ProxyResponseBuilderForCacheHit(r.Context(), gc.Answer, h.AuthDB, initTime, tokenLength, responseTime, apiKey)
	if err != nil {
		logrus.Errorf("Error building proxy response: %v", err)
		return
	}

	h.ChConfig.ResC <- &pRes

}
