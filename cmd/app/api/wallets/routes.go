package wallets

import (
	"github.com/go-chi/chi/v5"
)

func WalletsRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/add-funds", addFunds)
	r.Post("/withdraw-funds", withrawFunds)
	r.Get("/{user_id}", getWallet)
	r.Get("/{wallt_id}/history", getWalletHistory)
	r.Post("/bank-accounts", addBankAccount)
	r.Delete("/bank-accounts", deleteBankAccount)
	r.Post("/bank-accounts", getBankAccounts)

	return r
}
