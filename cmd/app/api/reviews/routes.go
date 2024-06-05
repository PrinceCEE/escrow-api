package reviews

import (
	"github.com/go-chi/chi/v5"
	"github.com/princecee/escrow-api/config"
)

func ReviewsRouter(c config.IConfig) chi.Router {
	h := reviewHandler{c}
	r := chi.NewRouter()

	r.Get("/not-implemented", h.notImplemented)

	return r
}
