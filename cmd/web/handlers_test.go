package main

import (
	"net/http"
	"net/url"
	"testing"

	"snippetbox.haonguyen.tech/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get(t, "/ping")
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	tests := []struct {
		name    string
		urlPath string
		expCode int
		expBody string
	}{
		{
			name:    "Valid ID",
			urlPath: "/snippet/1",
			expCode: http.StatusOK,
			expBody: "Bar",
		},
		{
			name:    "Non-existent ID",
			urlPath: "/snippet/2",
			expCode: http.StatusNotFound,
		},
		{
			name:    "Negative ID",
			urlPath: "/snippet/-1",
			expCode: http.StatusNotFound,
		},
		{
			name:    "Decimal ID",
			urlPath: "/snippet/1.23",
			expCode: http.StatusNotFound,
		},
		{
			name:    "String ID",
			urlPath: "/snippet/foo",
			expCode: http.StatusNotFound,
		},
		{
			name:    "Empty ID",
			urlPath: "/snippet/",
			expCode: http.StatusNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statusCode, _, body := ts.get(t, tc.urlPath)
			assert.Equal(t, statusCode, tc.expCode)
			assert.StringContains(t, body, tc.expBody)
		})
	}
}

func TestSignUpView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	// Make a get request to get the sign up form
	statusCode, _, body := ts.get(t, "/user/signup")
	assert.Equal(t, statusCode, http.StatusOK)

	// We need to have handlers to extract the CSRF token, because we have nosurf middleware which require these token for state changing request such as POST,PUT etc...
	validCSRFToken := extractCSRFToken(t, body)
	t.Logf("CSRF Token is :%q", validCSRFToken)

	const (
		validName     = "Foo"
		validPassword = "validPa$$word"
		validEmail    = "foo@gmail.com"
		formTag       = `<form action="/user/signup" method="POST" novalidate>`
	)
	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		expCode      int
		formTag      string
	}{
		{
			name:         "Valid Create User Submit",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			expCode:      http.StatusSeeOther,
		},
		{
			name:         "Invalid Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "Invalid Token",
			expCode:      http.StatusBadRequest,
		},
		{
			name:         "Empty Name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			expCode:      http.StatusUnprocessableEntity,
			formTag:      formTag,
		},
		{
			name:         "Empty Email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			expCode:      http.StatusUnprocessableEntity,
			formTag:      formTag,
		},
		{
			name:         "Empty Password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			expCode:      http.StatusUnprocessableEntity,
			formTag:      formTag,
		},
		{
			name:         "Invalid Email",
			userName:     validName,
			userEmail:    "invalidEmail",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			expCode:      http.StatusUnprocessableEntity,
			formTag:      formTag,
		},
		{
			name:         "Short Password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			csrfToken:    validCSRFToken,
			expCode:      http.StatusUnprocessableEntity,
			formTag:      formTag,
		},
		{
			name:         "Duplicate Email",
			userName:     validName,
			userEmail:    "exist@email.com",
			userPassword: "pa$$",
			csrfToken:    validCSRFToken,
			expCode:      http.StatusUnprocessableEntity,
			formTag:      formTag,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tc.userName)
			form.Add("email", tc.userEmail)
			form.Add("password", tc.userPassword)
			form.Add("name", tc.userName)
			form.Add("csrf_token", tc.csrfToken)
			statusCode, _, body := ts.postForm(t, "/user/signup", form)
			assert.Equal(t, statusCode, tc.expCode)
			if tc.formTag != "" {
				assert.StringContains(t, body, tc.formTag)
			}
		})
	}
}
