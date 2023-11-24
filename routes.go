package main

import (
	"github.com/Bupher-Co/bupher-api/api/auth"
	"github.com/Bupher-Co/bupher-api/api/businesses"
	"github.com/Bupher-Co/bupher-api/api/customers"
	"github.com/Bupher-Co/bupher-api/api/notifications"
	"github.com/Bupher-Co/bupher-api/api/reports"
	"github.com/Bupher-Co/bupher-api/api/reviews"
	"github.com/Bupher-Co/bupher-api/api/transactions"
	"github.com/Bupher-Co/bupher-api/api/users"
	"github.com/Bupher-Co/bupher-api/api/wallets"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

type routeFunc func(*config.Config) chi.Router

type routeConfig struct {
	path string
	f    routeFunc
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
		r.Mount(v.path, v.f(c))
	}

	return r
}
