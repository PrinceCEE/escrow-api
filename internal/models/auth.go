package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type PasswordHistory struct {
	Password  []byte    `json:"password"`
	Timestamp time.Time `json:"timestamp"`
}

type Auth struct {
	UserID          uuid.UUID         `json:"user_id"`
	Password        []byte            `json:"password"`
	PasswordHistory []PasswordHistory `json:"password_history"`
	ModelMixin
}
