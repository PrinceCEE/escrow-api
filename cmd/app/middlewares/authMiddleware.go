package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/princecee/escrow-api/cmd/app/pkg/response"
	"github.com/princecee/escrow-api/config"
	"github.com/princecee/escrow-api/pkg/jwt"
	"github.com/princecee/escrow-api/pkg/utils"
)

func AuthMiddleware(c config.IConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := response.ApiResponse{}

			auth := r.Header.Get("Authorization")
			if auth == "" {
				resp.Message = "authentication token missing"
				response.SendErrorResponse(w, resp, http.StatusUnauthorized)
				return
			}

			token := strings.Split(auth, " ")[1]
			claims, err := jwt.VerifyToken(token)
			if err != nil {
				resp.Message = err.Error()
				response.SendErrorResponse(w, resp, http.StatusUnauthorized)
				return
			}

			user, err := c.GetUserRepository().GetById(claims.UserID, nil)
			if err != nil {
				switch {
				case errors.Is(err, pgx.ErrNoRows):
					resp.Message = "invalid token"
				default:
					resp.Message = err.Error()
				}
				response.SendErrorResponse(w, resp, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), utils.ContextKey{}, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
