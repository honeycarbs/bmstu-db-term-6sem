package teststorage

import (
	"goodisgood/internal/app/model"
	"goodisgood/internal/app/storage"
)

type Storage struct {
	accountRepository *AccountRepository
	userRepository    *UserRepository
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Account() storage.AccountRepository {
	if s.accountRepository != nil {
		return s.accountRepository
	}

	s.accountRepository = &AccountRepository{
		storage:  s,
		accounts: make(map[int]*model.Account),
	}

	return s.accountRepository
}

func (s *Storage) User() storage.UserRepository {
	return nil
}
