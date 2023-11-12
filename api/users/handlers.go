package users

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type usersHandler struct {
	c *config.Config
}

func (ah *usersHandler) updateAccount(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (ah *usersHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
