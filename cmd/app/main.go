package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

func main() {
	r := getRouter()

	port := config.Config.Env.PORT
	srv := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	fmt.Printf("server running on port %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
