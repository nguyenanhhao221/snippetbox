package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// List of template files
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}

	// Convert all template paths to absolute paths
	for i, file := range files {
		absFilePath, err := filepath.Abs(file)
		if err != nil {
			app.errorLog.Println("Error finding absolute path:", err)
			app.serverError(w, err)
			return
		}
		files[i] = absFilePath
	}

	// Parse the template files
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	// Execute the "base" template and write it to the response
	if err := ts.ExecuteTemplate(w, "base", nil); err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d....\n", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Creating new snippet...\n")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
