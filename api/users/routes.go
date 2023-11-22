package users

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func UsersRouter(c *config.Config) chi.Router {
	uh := usersHandler{c}
	r := chi.NewRouter()

	r.Get("/me", uh.getMe)
	r.Get("/{user_id}", uh.getUser)
	r.Put("/update-account", uh.updateAccount)
	r.Put("/change-password", uh.changePassword)

	return r
}
