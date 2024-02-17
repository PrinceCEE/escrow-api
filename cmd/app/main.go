package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

func main() {
	r := getRouter()

	srv := http.Server{
		Addr:    ":" + config.Cfg.Env.PORT,
		Handler: r,
	}

	fmt.Printf("server running on port %s\n", config.Cfg.Env.PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
