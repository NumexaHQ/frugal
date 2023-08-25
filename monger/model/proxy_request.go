package model

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	nxAuthDB "github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/NumexaHQ/monger/utils"
	log "github.com/sirupsen/logrus"
)

type ProxyRequest struct {
	RequestID        string `json:"request_id"`
	RequestTimestamp int64  `json:"request_timestamp"`
	SourceIp         string `json:"source_ip"`
	RequestMethod    string `json:"request_method"`
	RequestUrl       string `json:"request_url"`
	RequestHeaders   string `json:"request_headers"`
	RequestBody      string `json:"request_body"`
	Provider         string `json:"provider"`
	UserID           int32  `json:"user_id"`
	ProjectID        int32  `json:"project_id"`
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

func ProxyRequestBuilderForHTTPRequest(r *http.Request, rt time.Time, authDB nxAuthDB.DB, url, apiKey string) (ProxyRequest, error) {
	ctx := r.Context()

	// get requset idfrom context
	rid := GetRequestID(ctx)

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
		buf.ReadFrom(r.Body)
		body = buf.String()
	}

	pr := ProxyRequest{
		RequestID:        rid,
		RequestTimestamp: rt.Unix(),
		SourceIp:         r.RemoteAddr,
		RequestMethod:    r.Method,
		RequestUrl:       url,
		RequestHeaders:   string(hB),
		RequestBody:      body,
		Provider:         "openai",
	}

	err = pr.SetUserIdentifier(ctx, authDB, apiKey)
	if err != nil {
		log.Errorf("Error setting user identifier: %v", err)
		return ProxyRequest{}, err
	}

	return pr, nil
}
