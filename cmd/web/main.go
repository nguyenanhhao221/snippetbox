package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the
// web application
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Flag to control config of the application
	addr := flag.String("addr", ":4000", "HTTP Network address")
	flag.Parse()

	// Custom logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)

	if err := srv.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}
