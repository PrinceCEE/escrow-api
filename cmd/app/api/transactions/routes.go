package transactions

import (
	"github.com/go-chi/chi/v5"
	"github.com/princecee/escrow-api/cmd/app/middlewares"
	"github.com/princecee/escrow-api/config"
)

func TransactionsRouter(c config.IConfig) chi.Router {
	t := transactionHandler{c}
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(c))

		r.Post("/create", t.createTransaction)
		r.Put("/{transction_id}", t.updateTransaction)
		r.Get("/{transaction_id}", t.getTransaction)
		r.Get("/", t.getTransactions)
		r.Post("/pay", t.makePayment)
	})

	return r
}
