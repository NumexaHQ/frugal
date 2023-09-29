package routes

import (
	nxAuthDB "github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/NumexaHQ/vibe/handlers"
	nxClickhouse "github.com/NumexaHQ/vibe/pkg/db/clickhouse"
	"github.com/gofiber/fiber/v2"
)

// Setup initializes the routes
func Setup(app *fiber.App, chConfig nxClickhouse.ClickhouseConfig, authdb nxAuthDB.DB, jwtSigningKey string) {
	nxHandler := &handlers.Handler{
		ChConfig:      chConfig,
		AuthDB:        authdb,
		JWTSigningKey: jwtSigningKey,
	}

	// livenies probe
	app.All("/ping", nxHandler.Pong)

	// Handlers for the management API
	app.Get("/mng_request/:userID", nxHandler.AuthMiddleware, nxHandler.GetRequestByUserID)
	app.Get("/mng_response/:requestID", nxHandler.AuthMiddleware, nxHandler.GetResponseByRequestID)

	// Handlers for the analytics/Metrics API
	app.Get("/avg_latency/:projectID", nxHandler.AuthMiddleware, nxHandler.ComputeAvgResponseLatencyByProjectID)
	app.Get("/total_requests/:projectID", nxHandler.AuthMiddleware, nxHandler.GetTotalRequestsCountbyProjectID)
	app.Get("/latency/:requestID", nxHandler.AuthMiddleware, nxHandler.ComputeLatencyByRequestId)
	app.Get("/avg_prompt_tokens/:projectID", nxHandler.AuthMiddleware, nxHandler.ComputeAverageTokensByProjectID)
	app.Get("/unique_models/:projectID", nxHandler.AuthMiddleware, nxHandler.GetUniqueModelsCountByProjectID)
	app.Get("/user_requests_stats/:projectID", nxHandler.AuthMiddleware, nxHandler.GetUserRequestsStatsByProjectID)
	app.Post("/add_prompt_directory", nxHandler.AuthMiddleware, nxHandler.AddRequestToPromptDirectory)
	app.Get("/prompt_directory/:projectID", nxHandler.AuthMiddleware, nxHandler.GetRequestFromPromptDirectory)
	app.Put("/update_prompt_directory", nxHandler.AuthMiddleware, nxHandler.EditFieldOfRequestInPromptDirectory)
	app.Get("/usage/:projectID", nxHandler.AuthMiddleware, nxHandler.GetUsageByProjectID)
}
