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
		app.serverError(w, err)
		return
	}

	// Execute the "base" template and write it to the response
	if err := ts.ExecuteTemplate(w, "base", nil); err != nil {
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
	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	id, err := app.model.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Log the redirect URL
	app.infoLog.Printf("Redirecting to /snippet/%d", id)
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
