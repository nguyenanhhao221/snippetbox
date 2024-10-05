package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// List of template files
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// Convert all template paths to absolute paths
	for i, file := range files {
		absFilePath, err := filepath.Abs(file)
		if err != nil {
			log.Println("Error finding absolute path:", err)
			http.Error(w, "Internal Server Error, Unable to Find Template Path", http.StatusInternalServerError)
			return
		}
		files[i] = absFilePath
	}

	// Parse the template files
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error, Fail To Parse HTML", http.StatusInternalServerError)
		return
	}

	// Execute the "base" template and write it to the response
	if err := ts.ExecuteTemplate(w, "base", nil); err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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
