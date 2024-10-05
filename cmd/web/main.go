package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Flag to control config of the application
	addr := flag.String("addr", ":4000", "HTTP Network address")
	flag.Parse()

	// Custom logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// as good practice always use our own ServerMux
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc(fmt.Sprintf("%s /snippet/{id}", http.MethodGet), snippetView)
	mux.HandleFunc(fmt.Sprintf("%s /snippet", http.MethodPost), snippetCreate)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s\n", *addr)

	if err := srv.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}
