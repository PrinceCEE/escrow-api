package users

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type usersHandler struct {
	c *config.Config
}

func (uh *usersHandler) getMe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (uh *usersHandler) getUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (uh *usersHandler) updateAccount(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (uh *usersHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
