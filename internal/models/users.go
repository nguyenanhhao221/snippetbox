package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// User Present the User in DB
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword string
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) error {
	hash_password, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, UTC_TIMESTAMP())"
	_, err = u.DB.Exec(stmt, name, email, hash_password)
	if err != nil {
		// Handle error, because email is UNIQUE Constrains , if duplicate email is found when insert, MYSQL will have a unique ERROR so we need to handle that case
		var mySqlError *mysql.MySQLError
		if errors.As(err, &mySqlError) {
			if mySqlError.Number == 1062 && strings.Contains(mySqlError.Message, "user_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate Check if the provided email and password match the credential in DB.
func (u *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword string
	stmt := "SELECT id, hashed_password  FROM users WHERE email = ?"
	err := u.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}
