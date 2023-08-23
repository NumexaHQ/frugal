package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/NumexaHQ/monger/model"

	"github.com/sirupsen/logrus"
)

// IngestLogs
func (h *Handler) IngestLogs(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("Error reading body: %v", err)
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	apiKey := r.Header.Get("X-Numexa-Api-Key")

	// read header if it is request or response
	// if request, build proxy request
	// if response, build proxy response
	if r.Header.Get("X-Numexa-Log-Type") == "request" {
		var embededRequest *http.Request
		err = json.Unmarshal(body, &embededRequest)
		if err != nil {
			logrus.Errorf("Error unmarshalling body: %v", err)
			http.Error(w, "Error unmarshalling body", http.StatusBadRequest)
			return
		}

		pr, err := model.ProxyRequestBuilderForHTTPRequest(embededRequest, time.Now(), h.AuthDB, r.URL.String(), apiKey)
		if err != nil {
			logrus.Errorf("Error building proxy request: %v", err)
		}

		go func() {
			h.ChConfig.ReqC <- &pr
		}()
		return
	} else if r.Header.Get("X-Numexa-Log-Type") == "response" {
		var embededResponse *http.Response
		err = json.Unmarshal(body, &embededResponse)
		if err != nil {
			logrus.Errorf("Error unmarshalling body: %v", err)
			http.Error(w, "Error unmarshalling body", http.StatusBadRequest)
			return
		}

		pr, err := model.ProxyRequestBuilderForHTTPResponse(r.Context(), embededResponse, h.AuthDB, time.Now(), 0, time.Now(), apiKey)
		if err != nil {
			logrus.Errorf("Error building proxy response: %v", err)
		}

		go func() {
			h.ChConfig.ResC <- &pr
		}()
		return
	}

	logrus.Errorf("Invalid log type")
	http.Error(w, "Invalid log type", http.StatusBadRequest)
}
