package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg"
	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var cfg = config.NewConfig()

func main() {
	apiRouter := initRoutes()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Mount("/api/v1", apiRouter)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		pkg.SendErrorResponse(w, pkg.ApiResponse{Message: fmt.Sprintf("%s %s not found", r.Method, r.URL.Path)}, http.StatusNotFound)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		pkg.SendErrorResponse(w, pkg.ApiResponse{Message: fmt.Sprintf("%s %s not allowed", r.Method, r.URL.Path)}, http.StatusMethodNotAllowed)
	})

	srv := http.Server{
		Addr:    ":" + cfg.Env.PORT,
		Handler: r,
	}

	fmt.Printf("server running on port %s\n", cfg.Env.PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
