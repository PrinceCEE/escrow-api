package jwt

import (
	"time"

	"github.com/Bupher-Co/bupher-api/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID    string `json:"user_id,omitempty"`
	Email     string `json:"email,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(t *TokenClaims) (string, error) {
	key := config.Config.Env.JWT_KEY

	t.IssuedAt = jwt.NewNumericDate(time.Now())
	t.Issuer = "bupherco"
	t.Subject = t.UserID

	if t.TokenType == "access_token" {
		t.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
	} else {
		t.ExpiresAt = jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	return token.SignedString(key)
}

func VerifyToken(tokenStr string) (*TokenClaims, error) {
	key := config.Config.Env.JWT_KEY

	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(*TokenClaims), nil
}
