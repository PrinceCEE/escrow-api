package businesses

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func BusinessRouter(c *config.Config) chi.Router {
	bh := businessHandler{c}
	r := chi.NewRouter()

	r.Get("/not-implemented", bh.NotImplemented)

	return r
}
