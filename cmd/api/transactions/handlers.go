package transactions

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
	"github.com/Bupher-Co/bupher-api/utils"
)

type transactionsHandler struct {
	c *config.Config
}

func (th *transactionsHandler) createTransaction(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func (th *transactionsHandler) updateTransaction(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func (th *transactionsHandler) acceptTransaction(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func (th *transactionsHandler) rejectTransaction(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func (th *transactionsHandler) makePayment(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func (th *transactionsHandler) getTransaction(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}

func (th *transactionsHandler) getTransactions(w http.ResponseWriter, r *http.Request) {
	utils.SendErrorResponse(w, utils.ApiResponse{Message: "not implemented"}, http.StatusNotImplemented)
}
