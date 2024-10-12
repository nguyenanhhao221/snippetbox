package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form"
)

// serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	if err := app.errorLog.Output(2, trace); err != nil {
		app.errorLog.Println("Fail to skip frame when log, fallback to full stack")
		app.errorLog.Println(trace)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		app.serverError(w, fmt.Errorf("Template %s does not exist", page))
		return
	}

	w.WriteHeader(status)
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) newTemplateData() *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// decodePostForm helper to decode form to dst, if dst is a nil pointer, panic
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError
		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}
		return err
	}
	return nil
}
