package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Bupher-Co/bupher-api/config"
)

func main() {
	c := config.NewConfig()
	defer c.DB.Close()

	r := getRouter(c)

	srv := http.Server{
		Addr:    ":" + c.Env.PORT,
		Handler: r,
	}

	fmt.Printf("server running on port %s\n", c.Env.PORT)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
