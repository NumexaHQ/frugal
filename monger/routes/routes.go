package routes

import (
	nxAuthDB "github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/NumexaHQ/monger/handlers"
	nxClickhouse "github.com/NumexaHQ/monger/pkg/db/clickhouse"
	"github.com/gorilla/mux"
)

// Setup initializes the routes
func Setup(r *mux.Router, chConfig nxClickhouse.ClickhouseConfig, authdb nxAuthDB.DB) {
	nxHandler := &handlers.Handler{
		ChConfig: chConfig,
		AuthDB:   authdb,
	}

	r.Use(Middleware)

	r.HandleFunc("/v1/openai/{rest:.*}", nxHandler.OpenAIProxy)

	// rest endpoint to ingest logs
	r.HandleFunc("/v1/logs", nxHandler.IngestLogs).Methods("POST")
}
