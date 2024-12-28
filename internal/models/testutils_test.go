package models

import (
	"database/sql"
	"os"
	"testing"
)

// newTestDB Setup and teardown for a test database to use for our test.
// Require mysql to run with user, table, database setup correctly
func newTestDB(t *testing.T) *sql.DB {
	db_conn_str := "test_web@/test_snippetbox?parseTime=true&multiStatements=true"
	db, err := sql.Open("mysql", db_conn_str)
	if err != nil {
		t.Fatal(err)
	}
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		script, err = os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})
	return db
}
