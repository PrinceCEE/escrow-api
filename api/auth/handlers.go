package auth

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type AuthHandler struct {
	c *config.Config
}

func NewAuthHandler(c *config.Config) *AuthHandler {
	return &AuthHandler{c}
}

func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (ah *AuthHandler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (ah *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (ah *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
