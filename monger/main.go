package main

import (
	"fmt"
	"net/http"
	"os"

	commonConstants "github.com/NumexaHQ/captainCache/numexa-common/constants"
	nxAuthDB "github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/NumexaHQ/monger/model"
	nxClickhouse "github.com/NumexaHQ/monger/pkg/db/clickhouse"
	"github.com/NumexaHQ/monger/pkg/worker"
	"github.com/NumexaHQ/monger/routes"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
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

	// Define a buffered channel with an appropriate buffer size
	var channelBufferSize = 100
	var reqC = make(chan *model.ProxyRequest, channelBufferSize)
	var resC = make(chan *model.ProxyResponse, channelBufferSize)

	chConfig := nxClickhouse.ClickhouseConfig{
		ReqC:     reqC,
		ResC:     resC,
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

	chConfig.DB = db

	chConfig.DB.AutoMigrate(&model.ProxyRequest{})
	chConfig.DB.AutoMigrate(&model.ProxyResponse{})

	// setup database connection to postgres
	psqldb := nxAuthDB.New(commonConstants.DBPostgres)
	if db == nil {
		log.Fatal("Failed to connect to psql database")
	}

	err = psqldb.Init()
	if err != nil {
		log.Fatal("Failed to initialize psql database")
	}

	// Start the worker
	w := worker.Worker{
		ChConfig: chConfig,
		AuthDB:   psqldb,
	}

	go w.Start()

	r := mux.NewRouter()

	// Initialize routes
	routes.Setup(r, chConfig, psqldb)

	// Start the server
	log.Fatal(http.ListenAndServe(":8081", r))
}
