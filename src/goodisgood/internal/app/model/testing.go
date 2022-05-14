package model

func TestAccount() *Account {
	return &Account{
		UUID:     "29a48ee1-6ba5-4c11-a208-b055bc1879a1",
		Username: "Gopher",
		Email:    "gopher@gopher.go",
		Password: "password",
	}
}

func TestUser() *User {
	return &User{
		UUID:   "9a81034c-2d81-42bf-8243-903ab73f5a91",
		Age:    20,
		Race:   "white",
		Gender: "M",
	}
}
