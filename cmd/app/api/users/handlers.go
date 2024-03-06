package users

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/response"
	"github.com/Bupher-Co/bupher-api/config"
)

type userHandler struct {
	c config.IConfig
}

func (h *userHandler) getMe(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *userHandler) getUser(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *userHandler) updateAccount(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *userHandler) changePassword(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}
