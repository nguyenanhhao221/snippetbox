package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"snippetbox.haonguyen.tech/internal/models"
)

// Define an application struct to hold the application-wide dependencies for the
// web application
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       models.SnippetModelInterface
	user           *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db_conn_str := os.Getenv("DB_CONNECTION")
	// Flag to control config of the application
	addr := flag.String("addr", ":4000", "HTTP Network address")
	// The parseTime=true part of the DSN above is a driver-specific parameter which instructs our driver to convert SQL TIME and DATE fields to Go time.Time objects.
	dns := flag.String("dns", db_conn_str, "MySQL connection")
	flag.Parse()

	// Custom logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dns)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	// Initialize session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Make sure that the Secure attribute is set on our session cookies.
	// Setting this means that the cookie will only be sent by a user's web
	// browser when a HTTPS connection is being used (and won't be sent over an
	// unsecure HTTP connection).
	sessionManager.Cookie.Secure = true

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		user: &models.UserModel{
			DB: db,
		},
		snippets: &models.SnippetModel{
			DB: db,
		},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Custom tls config for our server
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute, // All keep-alive connection will automatically drop after 1 minutes
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s\n", *addr)

	if err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"); err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(dns string) (*sql.DB, error) {
	// db here is a pool of connection
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
