package model_test

import (
	"goodisgood/internal/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		a       func() *model.Account
		isValid bool
	}{
		{
			name: "valid",
			a: func() *model.Account {
				return model.TestAccount()
			},
			isValid: true,
		},
		{
			name: "emptyEmail",
			a: func() *model.Account {
				a := model.TestAccount()
				a.Email = ""

				return a
			},
			isValid: false,
		},
		{
			name: "invalidEmail",
			a: func() *model.Account {
				a := model.TestAccount()
				a.Email = "jopa\n"

				return a
			},
			isValid: false,
		},
		{
			name: "emptyPassword",
			a: func() *model.Account {
				a := model.TestAccount()
				a.Password = ""

				return a
			},
			isValid: false,
		},
		{
			name: "shortPassword",
			a: func() *model.Account {
				a := model.TestAccount()
				a.Password = "111"

				return a
			},
			isValid: false,
		},
		{
			name: "withEncryptedPassword",
			a: func() *model.Account {
				a := model.TestAccount()
				a.Password = ""
				a.EncryptedPassword = "cant_steal_data_suck_it"

				return a
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.a().Validate())
			} else {
				assert.Error(t, tc.a().Validate())
			}
		})
	}
}

func TestEncryptPassword(t *testing.T) {
	a := model.TestAccount()
	assert.NoError(t, a.EncryptPassword())
	assert.NotEmpty(t, a.EncryptedPassword)
}
