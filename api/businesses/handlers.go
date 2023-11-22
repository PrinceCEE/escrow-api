package businesses

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type businessHandler struct {
	c *config.Config
}

func (bh *businessHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
