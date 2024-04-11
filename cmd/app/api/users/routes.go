package users

import (
	"github.com/Bupher-Co/bupher-api/cmd/app/middlewares"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func UsersRouter(c config.IConfig) chi.Router {
	h := userHandler{c}
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(c))

		r.Get("/me", h.getMe)
		r.Get("/{user_id}", h.getUser)
		r.Put("/update-account", h.updateAccount)
		r.Put("/change-password", h.changePassword)
	})

	return r
}
