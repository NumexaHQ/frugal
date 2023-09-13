package model

import (
	"time"
)

type User struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	OrganizationID int32     `json:"organization_id"`
	Email          string    `json:"email" validate:"required,email"`
	Password       string    `json:"password" validate:"required,min=8"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
