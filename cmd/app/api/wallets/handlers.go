package wallets

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/response"
)

func addBankAccount(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func deleteBankAccount(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func getBankAccounts(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func addFunds(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func withrawFunds(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func getWallet(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func getWalletHistory(w http.ResponseWriter, r *http.Request) {
	response.SendErrorResponse(w, response.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}
