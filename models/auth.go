package models

import "time"

type PasswordHistory struct {
	Password  string
	Timestamp time.Time
}

type Auth struct {
	UserID          string
	Password        string
	PasswordHistory []PasswordHistory
	ModelMixin
}
