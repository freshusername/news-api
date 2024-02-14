package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/freshusername/news-api/database"
)

const port = 3000

type application struct {
	DSN string
	DB  database.DatabaseRepo
}

func main() {
	// set application config
	var app application

	// read from command line
	dbName := os.Getenv("DB_NAME")
	flag.StringVar(&app.DSN, "dsn", fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable timezone=UTC connect_timeout=5",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		dbName,
		os.Getenv("DB_PORT"),
	), "Postgres connection string")
	flag.Parse()

	// connect to the database
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &database.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	log.Println("Starting application on port", port)

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
