package main

import (
	"context"
	"os"

	"github.com/NumexaHQ/captainCache/model"
	commonConstants "github.com/NumexaHQ/captainCache/numexa-common/constants"
	nxdb "github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/NumexaHQ/captainCache/routes"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
)

var redisClient *redis.Client

const (
	defaultJWTSigningKey = "change-me"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	// read jwt signing key from environment
	jwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	if jwtSigningKey == "" {
		jwtSigningKey = defaultJWTSigningKey
	}

	// init database
	db := nxdb.New(commonConstants.DBPostgres)
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	err := db.Init()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}

	// init AES setting
	err = model.InitializeAESSetting(context.Background(), db)
	if err != nil {
		log.Fatal("Failed to initialize AES setting: ", err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Use logger middleware
	app.Use(logger.New())

	// Set up CORS middleware with the correct AllowOrigins value
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Authorization,Content-Type",
		AllowCredentials: true,
	}))

	// Initialize routes
	routes.Setup(app, db, jwtSigningKey)

	// Start the server
	log.Fatal(app.Listen(":8080"))
}
