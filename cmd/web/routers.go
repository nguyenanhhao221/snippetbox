package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	// as good practice always use our own ServerMux
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Any route behind this will require authentication
	protected := dynamic.Append(app.requireAuth)

	mux.Handle("/", dynamic.ThenFunc(app.home))
	mux.Handle(fmt.Sprintf("%s /snippet/{id}", http.MethodGet), dynamic.ThenFunc(app.snippetView))
	mux.Handle(fmt.Sprintf("%s /snippet/latest", http.MethodGet), dynamic.ThenFunc(app.home))

	mux.Handle(fmt.Sprintf("%s /snippet", http.MethodGet), protected.ThenFunc(app.snippetCreateForm))
	mux.Handle(fmt.Sprintf("%s /snippet", http.MethodPost), protected.ThenFunc(app.snippetCreate))

	// Users
	mux.Handle(fmt.Sprintf("%s /user/signup", http.MethodGet), dynamic.ThenFunc(app.userSignUp))
	mux.Handle(fmt.Sprintf("%s /user/signup", http.MethodPost), dynamic.ThenFunc(app.userSignUpPost))
	mux.Handle(fmt.Sprintf("%s /user/login", http.MethodGet), dynamic.ThenFunc(app.userLogin))
	mux.Handle(fmt.Sprintf("%s /user/login", http.MethodPost), dynamic.ThenFunc(app.userLoginPost))

	mux.Handle(fmt.Sprintf("%s /user/logout", http.MethodPost), protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.panicRecover, app.loggingRequest, securityHeader)
	// We can apply middleware here by chaining
	return standard.Then(mux)
}
