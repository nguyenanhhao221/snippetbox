package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
	"snippetbox.haonguyen.tech/ui"
)

func (app *application) routes() http.Handler {
	// as good practice always use our own ServerMux
	mux := http.NewServeMux()

	// Handle static files by embed
	fileServer := http.FileServerFS(ui.Files)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("/static/", fileServer)
	// Ping route which does not request any middleware yet
	mux.HandleFunc(fmt.Sprintf("%s /ping", http.MethodGet), ping)

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
