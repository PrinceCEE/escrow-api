package models

import (
	"github.com/gofrs/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

type Token struct {
	Hash      []byte    `json:"hash" db:"hash"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TokenType TokenType `json:"token_type" db:"token_type"`
	InUse     bool      `json:"in_use" db:"in_use"`
	ModelMixin
}
