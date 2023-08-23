package handlers

import "github.com/NumexaHQ/captainCache/pkg/db"

type Handler struct {
	DB            db.DB
	JWTSigningKey string
}
