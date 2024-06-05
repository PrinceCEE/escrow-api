package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/princecee/escrow-api/cmd/app/middlewares"
	"github.com/princecee/escrow-api/config"
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
