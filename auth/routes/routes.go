// routes/routes.go
package routes

import (
	"github.com/NumexaHQ/captainCache/handlers"
	"github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Setup initializes the routes
func Setup(app *fiber.App, db db.DB, jwtSigningKey string) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	nxHandler := &handlers.Handler{
		DB:            db,
		JWTSigningKey: jwtSigningKey,
		Validator:     validate,
	}
	// Register and login handlers
	app.Post("/register", nxHandler.RegisterHandler)
	app.Post("/login", nxHandler.LoginHandler)
	app.Get("/logout", nxHandler.LogoutHandler)

	// Google OAuth handlers
	// app.Post("/google_auth", nxHandler.GoogleAuthCallback)

	//GenerateApiKey
	app.Post("/generate_api_key", nxHandler.AuthMiddleware, nxHandler.CreateApiKey)
	app.Get("/get_api_key", nxHandler.AuthMiddleware, nxHandler.GetAPIkeyByUserId)
}
