package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
}

func (app *application) serve() error {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Printf("Starting HTTP server in %s mode on port %d", app.config.env, app.config.port)

	return server.ListenAndServe()
}

func main() {
	var config config

	flag.IntVar(&config.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&config.env, "env", "development", "App enviroment {development|production|maintenance}")
	flag.StringVar(&config.db.dsn, "dsn", "", "Db url")

	flag.Parse()

	config.stripe.key = os.Getenv("STRIPE_KEY")
	config.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config:   config,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
	}

	err := app.serve()

	if err != nil {
		app.errorLog.Print(err)
		log.Fatal(err)
	}
}
