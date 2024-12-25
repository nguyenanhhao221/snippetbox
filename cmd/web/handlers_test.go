package main

import (
	"net/http"
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
