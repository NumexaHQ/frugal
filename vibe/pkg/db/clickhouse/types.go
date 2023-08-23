package clickhouse

import (
	"gorm.io/gorm"
)

type ClickhouseConfig struct {
	DB       *gorm.DB
	Username string
	Password string
	Hostname string
	Port     string
	Database string
	Address  string
}
