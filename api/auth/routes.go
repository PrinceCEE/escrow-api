package auth

import (
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func AuthRouter(c *config.Config) chi.Router {
	ah := AuthHandler{c}
	r := chi.NewRouter()

	r.Post("/sign-up", ah.SignUp)
	r.Post("/sign-in", ah.SignIn)
	r.Post("/verify-code", ah.VerifyCode)
	r.Post("/reset-password", ah.ResetPassword)
	r.Post("/change-password", ah.ChangePassword)

	return r
}
