package handlers

import (
	"github.com/NumexaHQ/captainCache/pkg/db"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	DB            db.DB
	JWTSigningKey string
	Validator     *validator.Validate
}
