package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.haonguyen.tech/internal/models"
	"snippetbox.haonguyen.tech/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	var data = app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 1,
	}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	form := snippetCreateForm{}
	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the embed Validator to validate form input
	form.CheckField(form.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(form.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(form.MaxChar(form.Title, 100), "title", "This field cannot more than 100 characters long")
	form.CheckField(form.PermittedInt(form.Expires, 7, 1, 365), "expires", "This field must equal 1 , 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Update the session to have a flash key
	app.sessionManager.Put(r.Context(), "flash", "Snippet created successfully!")
	// Log the redirect URL
	app.infoLog.Printf("Redirecting to /snippet/%d", id)
	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
