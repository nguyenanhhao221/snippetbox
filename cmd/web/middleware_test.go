package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.haonguyen.tech/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	mock := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			t.Fatal(err)
		}
	})
	securityHeader(mock).ServeHTTP(rr, r)

	rs := rr.Result()
	expValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expValue)

	expValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expValue)

	expValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expValue)

	expValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expValue)

	expValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expValue)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
