package reviews

import (
	"github.com/go-chi/chi/v5"
)

func ReviewsRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/not-implemented", notImplemented)

	return r
}
