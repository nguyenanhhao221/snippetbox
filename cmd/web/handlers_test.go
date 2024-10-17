package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.haonguyen.tech/internal/assert"
)

func TestPing(t *testing.T) {
	recorder := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(recorder, r)

	res := recorder.Result()
	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, string(body), "OK")
}
