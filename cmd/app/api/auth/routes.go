package auth

import (
	"github.com/go-chi/chi/v5"
)

func AuthRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/sign-up", signUp)
	r.Post("/sign-in", signIn)
	r.Post("/verify-code", verifyCode)
	r.Post("/reset-password", resetPassword)
	r.Post("/change-password", changePassword)

	return r
}
