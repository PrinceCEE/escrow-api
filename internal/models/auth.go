package models

import "time"

type PasswordHistory struct {
	Password  []byte    `json:"password"`
	Timestamp time.Time `json:"timestamp"`
}

type Auth struct {
	UserID          string            `json:"user_id"`
	Password        []byte            `json:"password"`
	PasswordHistory []PasswordHistory `json:"password_history"`
	ModelMixin
}
