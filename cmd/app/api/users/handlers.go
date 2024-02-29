package users

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/response"
)

func getMe(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}
