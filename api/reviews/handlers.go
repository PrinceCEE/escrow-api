package reviews

import (
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

type reviewsHandler struct {
	c *config.Config
}

func (rh *reviewsHandler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
