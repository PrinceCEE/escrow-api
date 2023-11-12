package users

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type UsersHandler struct {
	c *config.Config
}

func NewUsersHandler(c *config.Config) *UsersHandler {
	return &UsersHandler{c}
}

func (ah *UsersHandler) updateAccount(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (ah *UsersHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
