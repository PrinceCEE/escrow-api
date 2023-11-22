package transactions

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type transactionsHandler struct {
	c *config.Config
}

func (th *transactionsHandler) createTransaction(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (th *transactionsHandler) updateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (th *transactionsHandler) acceptTransaction(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (th *transactionsHandler) rejectTransaction(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (th *transactionsHandler) makePayment(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (th *transactionsHandler) getTransaction(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (th *transactionsHandler) getTransactions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
