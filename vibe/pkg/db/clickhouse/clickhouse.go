package clickhouse

import (
	"database/sql"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
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
