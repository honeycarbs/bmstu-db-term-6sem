package teststorage

import (
	"errors"
	"fmt"
	"goodisgood/internal/app/model"
	"strconv"
)

type AccountRepository struct {
	storage  *Storage
	accounts map[int]*model.Account
}

type UserRepository struct {
	storage *Storage
	users   map[int]*model.User
}

func (r *AccountRepository) Create(u *model.Account) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.EncryptPassword(); err != nil {
		return err
	}

	u.UUID = fmt.Sprint(len(r.accounts) + 1)
	r.accounts[len(r.accounts)+1] = u

	return nil
}

// func (r *AccountRepository) Find(id int) (*model.Account, error) {
// 	u, ok := r.accounts[id]
// 	if !ok {
// 		return nil, errors.New("not found")
// 	}

// 	return u, nil
// }

func (r *AccountRepository) FindByEmail(email string) (*model.Account, error) {
	for _, u := range r.accounts {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, errors.New("not found")
}

func (r *AccountRepository) Find(id string) (*model.Account, error) {
	id_, _ := strconv.Atoi(id)
	u, ok := r.accounts[id_]
	if !ok {
		return nil, errors.New("not found")
	}

	return u, nil
}
