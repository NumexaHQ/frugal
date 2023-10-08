package handlers

import (
	"github.com/NumexaHQ/monger/model"
	"io/ioutil"
	"net/http"

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
		err, newReq, requestTime := model.CreateNewRequest(err, body, r)
		if err != nil {
			logrus.Errorf("Error creating New Request: %v", err)
			http.Error(w, "Error creating New Request", http.StatusInternalServerError)
			return
		}
		pr, err := model.ProxyRequestBuilderForHTTPRequest(newReq, requestTime, h.AuthDB, newReq.URL.String(), apiKey)
		if err != nil {
			logrus.Errorf("Error building proxy request: %v", err)
		}

		go func() {
			h.ChConfig.ReqC <- &pr
		}()
		return
	} else if r.Header.Get("X-Numexa-Log-Type") == "response" {
		err, newResponse, initiatedTime, responseTime := model.CreateNewResponse(r, body)
		if err != nil {
			logrus.Errorf("Error building New Response: %v", err)
			http.Error(w, "Error building New Response", http.StatusInternalServerError)
			return
		}
		pr, err := model.ProxyResponseBuilderForHTTPResponse(r.Context(), newResponse, h.AuthDB, initiatedTime, 0, responseTime, apiKey)
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
