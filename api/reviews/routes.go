package reviews

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func ReviewsRouter(c *config.Config) chi.Router {
	rh := reviewsHandler{c}
	r := chi.NewRouter()

	r.Get("/not-implemented", rh.NotImplemented)

	return r
}
