package auth

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func AuthRouter(c *config.Config) chi.Router {
	ah := authHandler{c}
	r := chi.NewRouter()

	r.Post("/sign-up", ah.signUp)
	r.Post("/sign-in", ah.signIn)
	r.Post("/verify-code", ah.verifyCode)
	r.Post("/reset-password", ah.resetPassword)
	r.Post("/change-password", ah.changePassword)

	return r
}
