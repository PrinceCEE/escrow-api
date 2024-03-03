package models

import (
	"github.com/gofrs/uuid"
)

type Business struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	ModelMixin
}
