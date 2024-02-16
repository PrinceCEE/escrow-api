package customers

import (
	"github.com/go-chi/chi/v5"
)

func CustomerRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/not-implemented", notImplemented)

	return r
}
