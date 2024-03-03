package customers

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func CustomerRouter(c *config.Config) chi.Router {
	h := customerHandler{c}
	r := chi.NewRouter()

	r.Get("/not-implemented", h.notImplemented)

	return r
}
