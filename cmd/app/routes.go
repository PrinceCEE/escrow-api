package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Bupher-Co/bupher-api/cmd/app/api/auth"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/businesses"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/customers"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/notifications"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/reports"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/reviews"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/transactions"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/users"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/wallets"
	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/response"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

type routeFunc func(*config.Config) chi.Router

type routeConfig struct {
	path string
	fn   routeFunc
}

func initRoutes(c *config.Config) chi.Router {
	routes := []routeConfig{
		{"/auth", auth.AuthRouter},
		{"/users", users.UsersRouter},
		{"/wallets", wallets.WalletsRouter},
		{"/notifications", notifications.NotificationRouter},
		{"/transactions", transactions.TransactionsRouter},
		{"/businesses", businesses.BusinessRouter},
		{"/customers", customers.CustomerRouter},
		{"/reviews", reviews.ReviewsRouter},
		{"/reports", reports.ReportRouter},
	}

	r := chi.NewRouter()
	for _, v := range routes {
		r.Mount(v.path, v.fn(c))
	}

	return r
}

func getRouter(c *config.Config) chi.Router {
	apiRouter := initRoutes(c)
	r := chi.NewRouter()

	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(middleware.CleanPath)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)
	r.Mount("/api/v1", apiRouter)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		resp := response.ApiResponse{Message: fmt.Sprintf("%s %s not found", r.Method, r.URL.Path)}
		response.SendErrorResponse(w, resp, http.StatusNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		resp := response.ApiResponse{Message: fmt.Sprintf("%s %s not allowed", r.Method, r.URL.Path)}
		response.SendErrorResponse(w, resp, http.StatusMethodNotAllowed)
	})

	return r
}
