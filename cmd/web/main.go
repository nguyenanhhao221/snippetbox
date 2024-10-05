package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// as good practice always use our own ServerMux
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc(fmt.Sprintf("%s /snippet/{id}", http.MethodGet), snippetView)
	mux.HandleFunc(fmt.Sprintf("%s /snippet", http.MethodPost), snippetCreate)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
