package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {
	
}

func main() {
	app := Config{}

	log.Println("Starting mail service on port: 80")

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
