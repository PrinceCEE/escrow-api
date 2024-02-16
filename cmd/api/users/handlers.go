package users

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/utils"
)

type usersHandler struct {
	c *config.Config
}

func (uh *usersHandler) getMe(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func (uh *usersHandler) getUser(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func (uh *usersHandler) updateAccount(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func (uh *usersHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}
