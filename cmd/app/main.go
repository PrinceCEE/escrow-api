package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

var cfg = config.NewConfig()

func main() {
	r := getRouter()

	srv := http.Server{
		Addr:    ":" + cfg.Env.PORT,
		Handler: r,
	}

	fmt.Printf("server running on port %s\n", cfg.Env.PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
