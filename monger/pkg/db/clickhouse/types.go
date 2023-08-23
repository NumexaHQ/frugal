package clickhouse

import (
	"github.com/NumexaHQ/monger/model"
	"gorm.io/gorm"
)

type ClickhouseConfig struct {
	DB       *gorm.DB
	ReqC     chan *model.ProxyRequest
	ResC     chan *model.ProxyResponse
	Username string
	Password string
	Hostname string
	Port     string
	Database string
	Address  string
}
