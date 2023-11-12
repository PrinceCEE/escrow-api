package main

import (
	"log"
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
	"github.com/go-chi/chi/v5"
)

func main() {
	c := config.NewConfig()

	apiRouter := initRoutes(c)
	r := chi.NewRouter()
	r.Mount("/api/v1", apiRouter)

	srv := http.Server{
		Addr: ":" + c.Env.PORT,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
