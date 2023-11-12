package wallets

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type walletsHandler struct {
	c *config.Config
}

func (wh *walletsHandler) addBankAccount(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (wh *walletsHandler) deleteBankAccount(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (wh *walletsHandler) getBankAccounts(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (wh *walletsHandler) addFunds(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (wh *walletsHandler) withrawFunds(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (wh *walletsHandler) getWallet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}

func (wh *walletsHandler) getWalletHistory(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
