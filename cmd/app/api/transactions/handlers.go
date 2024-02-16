package transactions

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg"
)

func createTransaction(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func updateTransaction(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func acceptTransaction(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func rejectTransaction(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func makePayment(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
	pkg.SendErrorResponse(w, pkg.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}
