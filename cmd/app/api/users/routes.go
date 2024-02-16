package users

import (
	"github.com/go-chi/chi/v5"
)

func UsersRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/me", getMe)
	r.Get("/{user_id}", getUser)
	r.Put("/update-account", updateAccount)
	r.Put("/change-password", changePassword)

	return r
}
