package model

import validation "github.com/go-ozzo/ozzo-validation"

type User struct {
	UUID   string `json:"uuid,omitempty"`
	Age    int    `json:"age"`
	Race   string `json:"race"`
	Gender string `json:"gender"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Gender, validation.In("F", "M")),
	)
}
