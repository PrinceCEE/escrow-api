package main

import (
	"github.com/Bupher-Co/bupher-api/cmd/app/api/auth"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/businesses"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/customers"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/notifications"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/reports"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/reviews"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/transactions"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/users"
	"github.com/Bupher-Co/bupher-api/cmd/app/api/wallets"
	"github.com/go-chi/chi/v5"
)

type routeFunc func() chi.Router

type routeConfig struct {
	path string
	f    routeFunc
}

func initRoutes() chi.Router {
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
		r.Mount(v.path, v.f())
	}

	return r
}
