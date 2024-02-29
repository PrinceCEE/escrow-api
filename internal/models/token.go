package models

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

type Token struct {
	Hash      []byte    `json:"hash"`
	UserID    string    `json:"user_id"`
	TokenType TokenType `json:"token_type"`
	InUse     bool      `json:"in_use"`
	ModelMixin
}
