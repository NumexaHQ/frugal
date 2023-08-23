package handlers

import (
	nxAuthDB "github.com/NumexaHQ/captainCache/pkg/db"
	nxClickhouse "github.com/NumexaHQ/vibe/pkg/db/clickhouse"
)

type Handler struct {
	ChConfig      nxClickhouse.ClickhouseConfig
	AuthDB        nxAuthDB.DB
	JWTSigningKey string
}
