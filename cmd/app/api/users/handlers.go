package users

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg"
)

func getMe(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}
