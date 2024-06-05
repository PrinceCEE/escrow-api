package customers

import (
	"github.com/go-chi/chi/v5"
	"github.com/princecee/escrow-api/config"
)

func CustomerRouter(c config.IConfig) chi.Router {
	h := customerHandler{c}
	r := chi.NewRouter()

	r.Get("/not-implemented", h.notImplemented)

	return r
}
