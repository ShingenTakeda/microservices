package main

import (
	"fmt"
	"log"
	"net/http"
)

const web_port = "80"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting broker service in port: %s\n", web_port)

	// Define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", web_port),
		Handler: app.routes(),
	}

	// Start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
