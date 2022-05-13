package model

func TestAccount() *Account {
	return &Account{
		UUID:     "29a48ee1-6ba5-4c11-a208-b055bc1879a1",
		Email:    "gopher@gopher.go",
		Password: "password",
	}
}
