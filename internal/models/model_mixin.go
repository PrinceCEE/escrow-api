package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type ModelMixin struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
	Version   int64     `json:"version" db:"version"`
}
