package wallets

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func WalletsRouter(c *config.Config) chi.Router {
	wh := walletsHandler{c}
	r := chi.NewRouter()

	r.Post("/add-funds", wh.addFunds)
	r.Post("withdraw-funds", wh.withrawFunds)
	r.Get("{user_id}", wh.getWallet)
	r.Get("{wallt_id}/history", wh.getWalletHistory)
	r.Post("/bank-accounts", wh.addBankAccount)
	r.Delete("/bank-accounts", wh.deleteBankAccount)
	r.Post("/bank-accounts", wh.getBankAccounts)

	return r
}
