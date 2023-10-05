package model

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	nxAuthDB "github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/NumexaHQ/monger/utils"
	log "github.com/sirupsen/logrus"
)

type ProxyRequest struct {
	RequestID        string    `json:"request_id"`
	RequestTimestamp time.Time `json:"request_timestamp"`
	SourceIp         string    `json:"source_ip"`
	RequestMethod    string    `json:"request_method"`
	RequestUrl       string    `json:"request_url"`
	RequestHeaders   string    `json:"request_headers"`
	RequestBody      string    `json:"request_body"`
	Provider         string    `json:"provider"`
	UserID           int32     `json:"user_id"`
	ProjectID        int32     `json:"project_id"`
	IsCached         bool      `json:"is_cached"`
	IsCacheHit       bool      `json:"is_cache_hit"`
	CustomMetaData   string    `json:"custom_fields"`
}

func (p *ProxyRequest) SetUserIdentifier(ctx context.Context, authDB nxAuthDB.DB, apiKey string) error {
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

func ProxyRequestBuilderForHTTPRequest(r *http.Request, rt time.Time, authDB nxAuthDB.DB, url, apiKey string, index int) (ProxyRequest, error) {
	ctx := r.Context()

	// get requset idfrom context
	rid := GetRequestID(ctx)

	//check for cache header
	isCached := false
	if r.Header.Get("X-Numexa-Cache") == "true" {
		isCached = true
	}

	// custom metadata from the header
	metaDataHeader := r.Header.Get("X-Numexa-Custom-Metadata")
	var metaData map[string]interface{}
	if metaDataHeader != "" {
		err := json.Unmarshal([]byte(metaDataHeader), &metaData)
		if err != nil {
			log.Errorf("Error unmarshalling metadata header: %v", err)
			return ProxyRequest{}, err
		}
	}

	// change metaDataHeader to JSON string
	metaDataHeaderB, err := json.Marshal(metaData)
	if err != nil {
		log.Errorf("Error marshalling metadata header: %v", err)
		return ProxyRequest{}, err
	}
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

	hB, err := json.Marshal(header)
	if err != nil {
		log.Errorf("Error marshalling header: %v", err)
		return ProxyRequest{}, err
	}

	body := ""
	if r.Body != nil {
		buf := new(bytes.Buffer)
		requestBodyBytes, _ := io.ReadAll(r.Body)
		requestBody := &utils.RequestBody{}
		if err := json.Unmarshal(requestBodyBytes, requestBody); err != nil {
			return ProxyRequest{}, err
		}
		newRequestBody := &utils.FinalRequestBody{
			Model:    requestBody.Config.Options[index].OverrideParams.Model,
			Messages: requestBody.Params.Messages,
		}
		newRequestBodyBytes, _ := json.Marshal(newRequestBody)
		_, _ = buf.ReadFrom(bytes.NewReader(newRequestBodyBytes))
		body = buf.String()
	}

	pr := ProxyRequest{
		RequestID:        rid,
		RequestTimestamp: rt,
		SourceIp:         r.RemoteAddr,
		RequestMethod:    r.Method,
		RequestUrl:       url,
		RequestHeaders:   string(hB),
		RequestBody:      body,
		Provider:         "openai",
		IsCached:         isCached,
		CustomMetaData:   string(metaDataHeaderB),
	}

	err = pr.SetUserIdentifier(ctx, authDB, apiKey)
	if err != nil {
		log.Errorf("Error setting user identifier: %v", err)
		return ProxyRequest{}, err
	}

	return pr, nil
}
