package wallets

import (
	"github.com/go-chi/chi/v5"
	"github.com/princecee/escrow-api/cmd/app/middlewares"
	"github.com/princecee/escrow-api/config"
)

func WalletsRouter(c config.IConfig) chi.Router {
	h := walletHandler{c}
	r := chi.NewRouter()

	r.Post("/paystack-webhook", h.handlePaystackWebhook)

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(c))

		r.Post("/add-funds", h.addFunds)
		r.Post("/withdraw-funds", h.withrawFunds)
		r.Get("/", h.getWallet)
		r.Get("/{wallet_id}/history", h.getWalletHistories)
		r.Post("/bank-accounts", h.addBankAccount)
		r.Delete("/bank-accounts/{bank_account_id}", h.deleteBankAccount)
		r.Get("/bank-accounts", h.getBankAccounts)
	})

	return r
}
