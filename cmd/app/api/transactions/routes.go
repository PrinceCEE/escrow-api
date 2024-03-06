package transactions

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func TransactionsRouter(c config.IConfig) chi.Router {
	h := transactionHandler{c}
	r := chi.NewRouter()

	r.Post("/create", h.createTransaction)
	r.Put("/{transction_id}", h.updateTransaction)
	r.Post("/{transaction_id}/accpet", h.acceptTransaction)
	r.Post("/{transaction_id}/reject", h.rejectTransaction)
	r.Get("/{transaction_id}", h.getTransaction)
	r.Get("/", h.getTransactions)
	r.Post("/pay", h.makePayment)

	return r
}
