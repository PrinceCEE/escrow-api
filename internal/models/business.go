package models

import (
	"github.com/gofrs/uuid"
)

type Business struct {
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Name   string    `json:"name" db:"name"`
	Email  string    `json:"email" db:"email"`
	ModelMixin
}
