package reviews

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func ReviewsRouter(c *config.Config) chi.Router {
	h := reviewHandler{c}
	r := chi.NewRouter()

	r.Get("/not-implemented", h.notImplemented)

	return r
}
