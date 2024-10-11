package main

import (
	"fmt"
	"net/http"
)

func (app *application) routes() http.Handler {

	// as good practice always use our own ServerMux
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc(fmt.Sprintf("%s /snippet/{id}", http.MethodGet), app.snippetView)
	mux.HandleFunc(fmt.Sprintf("%s /snippet/latest", http.MethodGet), app.snippetViewLatest)
	mux.HandleFunc(fmt.Sprintf("%s /snippet", http.MethodPost), app.snippetCreate)

	// Apply securityHeader middleware
	return securityHeader(mux)
}
