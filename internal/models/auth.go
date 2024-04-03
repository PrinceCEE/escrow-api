package models

import (
	"time"
)

type PasswordHistory struct {
	Password  string    `json:"password"`
	Timestamp time.Time `json:"timestamp"`
}

type Auth struct {
	UserID          string            `json:"user_id" db:"user_id"`
	Password        string            `json:"password" db:"password"`
	PasswordHistory []PasswordHistory `json:"password_history" db:"password_history"`
	ModelMixin
}
