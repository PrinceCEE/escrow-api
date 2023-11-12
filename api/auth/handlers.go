package auth

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type authHandler struct {
	c *config.Config
}

func (ah *authHandler) signUp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (ah *authHandler) signIn(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (ah *authHandler) verifyCode(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (ah *authHandler) resetPassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (ah *authHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
