package wallets

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/response"
	"github.com/Bupher-Co/bupher-api/config"
)

type walletHandler struct {
	c *config.Config
}

func (h *walletHandler) addBankAccount(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) deleteBankAccount(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) getBankAccounts(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) addFunds(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) withrawFunds(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) getWallet(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}

func (h *walletHandler) getWalletHistory(w http.ResponseWriter, r *http.Request) {
	resp := response.ApiResponse{Message: "not implemented"}
	response.SendErrorResponse(w, resp, http.StatusNotImplemented)
}
