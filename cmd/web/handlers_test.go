package main

import (
	"net/http"
	"testing"

	"snippetbox.haonguyen.tech/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication()

	ts := newTestServer(app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get(t, "/ping")
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")
}
