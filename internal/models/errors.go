package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrDuplicateEmail     = errors.New("models: email already exist")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
)
