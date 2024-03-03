package auth

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func AuthRouter(c *config.Config) chi.Router {
	h := authHandler{c}
	r := chi.NewRouter()

	r.Post("/sign-up", h.signUp)
	r.Post("/sign-in", h.signIn)
	r.Post("/verify-code", h.verifyCode)
	r.Post("/reset-password", h.resetPassword)
	r.Post("/change-password", h.changePassword)

	return r
}
