package customers

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type customersHandler struct {
	c *config.Config
}

func (ch *customersHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
