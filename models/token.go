package models

type TokenType string

const (
	AccessToken  TokenType = "access_token"
	RefreshToken TokenType = "refresh_token"
)

type Token struct {
	Hash      string    `json:"hash"`
	TokenType TokenType `json:"token_type"`
	UserID    string    `json:"user_id"`
	InUse     bool      `json:"in_use"`
	ModelMixin
}
