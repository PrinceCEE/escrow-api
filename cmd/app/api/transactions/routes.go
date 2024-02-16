package transactions

import (
	"github.com/go-chi/chi/v5"
)

func TransactionsRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/create", createTransaction)
	r.Put("/{transction_id}", updateTransaction)
	r.Post("/{transaction_id}/accpet", acceptTransaction)
	r.Post("/{transaction_id}/reject", rejectTransaction)
	r.Get("/{transaction_id}", getTransaction)
	r.Get("/", getTransactions)
	r.Post("/pay", makePayment)

	return r
}
