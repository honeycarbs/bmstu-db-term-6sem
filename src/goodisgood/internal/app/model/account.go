package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	UUID              string `json:"uuid"`
	Username          string `json:"username,omitempty"`
	Email             string `json:"email,omitempty"`
	Password          string `json:"password,omitempty"`
	Role              string `json:"role,omitempty"`
	EncryptedPassword string `json:"-"`
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

func (a *Account) Sanitize() {
	a.Password = ""
}

func (a *Account) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password)) == nil
}

func encryptString(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
