package clickhouse

import (
	"context"
	"database/sql"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/NumexaHQ/monger/model"
)

func (c *ClickhouseConfig) OpenDB() *sql.DB {
	sqlDB := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{c.Address},
		Auth: clickhouse.Auth{
			Database: c.Database,
			Username: c.Username,
			Password: c.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Debug:       true,
	})
	return sqlDB
}

func (c *ClickhouseConfig) IngestProxyRequest(ctx context.Context, request model.ProxyRequest) error {
	c.DB.Create(&request)
	return nil
}

func (c *ClickhouseConfig) IngestProxyResponse(ctx context.Context, response model.ProxyResponse) error {
	c.DB.Create(&response)
	return nil
}
