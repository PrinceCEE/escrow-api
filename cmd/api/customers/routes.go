package customers

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func CustomerRouter(c *config.Config) chi.Router {
	ch := customersHandler{c}
	r := chi.NewRouter()

	r.Get("/not-implemented", ch.NotImplemented)

	return r
}
