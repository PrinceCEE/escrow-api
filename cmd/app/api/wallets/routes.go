package wallets

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func WalletsRouter(c *config.Config) chi.Router {
	h := walletHandler{c}
	r := chi.NewRouter()

	r.Post("/add-funds", h.addFunds)
	r.Post("/withdraw-funds", h.withrawFunds)
	r.Get("/{user_id}", h.getWallet)
	r.Get("/{wallt_id}/history", h.getWalletHistory)
	r.Post("/bank-accounts", h.addBankAccount)
	r.Delete("/bank-accounts", h.deleteBankAccount)
	r.Post("/bank-accounts", h.getBankAccounts)

	return r
}
