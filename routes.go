package main

import (
	"github.com/Bupher-Co/bupher-api/api/auth"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func initRoutes(c *config.Config) chi.Router {
	r := chi.NewRouter()
	ar := auth.AuthRouter(c)

	// mount the routers
	r.Mount("/auth", ar)

	return r
}