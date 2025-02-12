package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const web_port = "80"

var retry_count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting auth service")

	conn := connect_to_DB()

	if conn == nil {
		log.Panic("Couldn't connect to Postgres")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", web_port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func open_DB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connect_to_DB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := open_DB(dsn)
		if err != nil {
			log.Println("Postgres not ready")
			retry_count++
		} else {
			log.Println("Connected to postgres")
			return connection
		}

		if retry_count > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Waiting for 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
