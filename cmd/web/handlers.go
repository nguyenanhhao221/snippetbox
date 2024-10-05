package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if _, err := w.Write([]byte("Hello from Snippetbox.\n")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d....\n", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Creating new snippet...\n")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
