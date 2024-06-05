package businesses

import (
	"github.com/go-chi/chi/v5"
	"github.com/princecee/escrow-api/config"
)

func BusinessRouter(c config.IConfig) chi.Router {
	h := businessHandler{c}
	r := chi.NewRouter()

	r.Get("/not-implemented", h.notImplemented)

	return r
}
