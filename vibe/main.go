package main

import (
	"fmt"
	"os"

	commonConstants "github.com/NumexaHQ/captainCache/numexa-common/constants"
	nxAuthDB "github.com/NumexaHQ/captainCache/pkg/db"
	nxClickhouse "github.com/NumexaHQ/vibe/pkg/db/clickhouse"
	"github.com/NumexaHQ/vibe/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

const (
	defaultJWTSigningKey = "change-me"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	// get the clickhouse details from the environment
	username := os.Getenv("CLICKHOUSE_USER")
	password := os.Getenv("CLICKHOUSE_PASSWORD")
	hostname := os.Getenv("CLICKHOUSE_HOST")
	port := os.Getenv("CLICKHOUSE_PORT")
	database := os.Getenv("CLICKHOUSE_DB")

	// username = "numexa"
	// password = "numexa"
	// hostname = "localhost"
	// port = "9000"
	// database = "numexa"

	// read jwt signing key from environment
	jwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	if jwtSigningKey == "" {
		jwtSigningKey = defaultJWTSigningKey
	}

	chConfig := nxClickhouse.ClickhouseConfig{
		Username: username,
		Password: password,
		Hostname: hostname,
		Port:     port,
		Database: database,
		Address:  fmt.Sprintf("%s:%s", hostname, port),
	}

	// Set up the database connection to clickhouse
	sqlDB := chConfig.OpenDB()

	db, err := gorm.Open(clickhouse.New(clickhouse.Config{
		Conn: sqlDB,
	}))
	if err != nil {
		log.Fatal(err)
	}
	// setup database connection to postgres
	psqldb := nxAuthDB.New(commonConstants.DBPostgres)
	if db == nil {
		log.Fatal("Failed to connect to psql database")
	}

	err = psqldb.Init()
	if err != nil {
		log.Fatal("Failed to initialize psql database")
	}

	chConfig.DB = db

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
	routes.Setup(app, chConfig, psqldb, jwtSigningKey)

	// Start the server
	log.Fatal(app.Listen(":8082"))
}
