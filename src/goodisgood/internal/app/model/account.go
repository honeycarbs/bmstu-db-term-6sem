package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	UUID              string
	Username          string
	Email             string
	Password          string
	EncryptedPassword string
}

func (a *Account) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.By(requiredIf(a.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

func (a *Account) EncryptPassword() error {
	if len(a.Password) > 0 {
		enc, err := encryptString(a.Password)
		if err != nil {
			return err
		}

		a.EncryptedPassword = enc
	}
	return nil
}

func encryptString(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
