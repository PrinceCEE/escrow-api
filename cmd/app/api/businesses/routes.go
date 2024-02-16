package businesses

import (
	"github.com/go-chi/chi/v5"
)

func BusinessRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/not-implemented", notImplemented)

	return r
}
