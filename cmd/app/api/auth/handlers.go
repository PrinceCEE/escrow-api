package auth

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg"
)

func signUp(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func verifyCode(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}
