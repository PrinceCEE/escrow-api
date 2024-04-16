package transactions

import (
	"github.com/Bupher-Co/bupher-api/cmd/app/middlewares"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
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
