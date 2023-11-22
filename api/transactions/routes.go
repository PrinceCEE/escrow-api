package transactions

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func TransactionsRouter(c *config.Config) chi.Router {
	th := transactionsHandler{c}
	r := chi.NewRouter()

	r.Post("/create", th.createTransaction)
	r.Put("/{transction_id}", th.updateTransaction)
	r.Post("/{transaction_id}/accpet", th.acceptTransaction)
	r.Post("/{transaction_id}/reject", th.rejectTransaction)
	r.Get("/{transaction_id}", th.getTransaction)
	r.Get("/", th.getTransactions)
	r.Post("/pay", th.makePayment)

	return r
}
