package worker

import (
	"context"
	"sync"

	authDB "github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/NumexaHQ/monger/model"
	nxClickhouse "github.com/NumexaHQ/monger/pkg/db/clickhouse"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	ChConfig  nxClickhouse.ClickhouseConfig
	AuthDB    authDB.DB
	waitGroup *sync.WaitGroup
}

func (w *Worker) Start() {
	ctx := context.Background()

	// wait group
	w.waitGroup = &sync.WaitGroup{}

	w.waitGroup.Add(2)

	go w.sendRequest(ctx, w.ChConfig.ReqC)
	go w.sendResponse(ctx, w.ChConfig.ResC)

	w.waitGroup.Wait()
	log.Info("Worker finished")
}

func (w *Worker) sendRequest(ctx context.Context, reqC chan *model.ProxyRequest) {
	log.Info("Starting to send requests to ClickHouse")
	for req := range reqC {
		// Insert the data into ClickHouse here
		log.Debugf("Sending request to ClickHouse: %v", req)
		w.ChConfig.IngestProxyRequest(ctx, *req)
	}
	defer w.waitGroup.Done()
}

func (w *Worker) sendResponse(ctx context.Context, resC chan *model.ProxyResponse) {
	log.Info("Starting to send responses to ClickHouse")
	for res := range resC {
		// Insert the data into ClickHouse here
		log.Debugf("Sending response to ClickHouse: %v", res)
		w.ChConfig.IngestProxyResponse(ctx, *res)
	}
	defer w.waitGroup.Done()
}
