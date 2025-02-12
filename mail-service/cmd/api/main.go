package main

import (
	"log"
	"net/http"
)

type Config struct{}

const web_port = "80"

func main() {
	app := Config{}

	log.Println("starting on port:", web_port)

	srv := &http.Server{
		Addr:    web_port,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
