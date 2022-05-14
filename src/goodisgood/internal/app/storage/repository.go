package storage

import "goodisgood/internal/app/model"

type AccountRepository interface {
	Create(*model.Account) error
	FindByEmail(string) (*model.Account, error)
	Find(string) (*model.Account, error)
}

type UserRepository interface {
	Create(string, *model.User) error
}
