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
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	})

	srv := http.Server{
		Addr:    ":" + c.Env.PORT,
		Handler: r,
	}

	log.Printf("server running on port %s\n", c.Env.PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
