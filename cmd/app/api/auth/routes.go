package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/princecee/escrow-api/config"
)

func AuthRouter(c config.IConfig) chi.Router {
	h := authHandler{c}
	r := chi.NewRouter()

	r.Post("/sign-up", h.signUp)
	r.Post("/sign-in", h.signIn)
	r.Post("/verify-code", h.verifyCode)
	r.Post("/reset-password", h.resetPassword)
	r.Post("/change-password", h.changePassword)
	r.Post("/resend-code", h.resendCode)

	return r
}
