package businesses

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func BusinessRouter(c *config.Config) chi.Router {
	h := businessHandler{c}
	r := chi.NewRouter()

	r.Get("/not-implemented", h.notImplemented)

	return r
}
