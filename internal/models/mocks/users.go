package mocks

import (
	"snippetbox.haonguyen.tech/internal/models"
)

type UserModel struct{}

func (u *UserModel) Insert(name, email, password string) error {
	switch email {
	case "exist@email.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	if email == "foo@gmail.com" && password == "validPa$$word" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredentials
}

func (u *UserModel) Exist(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
