package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Bupher-Co/bupher-api/cmd/app/pkg/routes"
	"github.com/Bupher-Co/bupher-api/config"
)

func main() {
	c := config.NewConfig()
	defer c.DB.Close()

	r := routes.GetRouter(c)

	port := c.Getenv("PORT")
	srv := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	fmt.Printf("server running on port %s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
