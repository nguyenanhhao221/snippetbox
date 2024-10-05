package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Flag to control config of the application
	addr := flag.String("addr", ":4000", "HTTP Network address")
	flag.Parse()

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

	log.Printf("Starting server on %s\n", *addr)
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
