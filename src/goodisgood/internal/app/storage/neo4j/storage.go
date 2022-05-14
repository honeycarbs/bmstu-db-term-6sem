package neo4jstorage

import (
	"goodisgood/internal/app/storage"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Storage struct {
	db                neo4j.Driver
	accountRepository *AccountRepository
	userRepository    *UserRepository
}

func NewStorage(db neo4j.Driver) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Account() storage.AccountRepository {
	if s.accountRepository != nil {
		return s.accountRepository
	}

	s.accountRepository = &AccountRepository{
		storage: s,
	}

	return s.accountRepository
}

func (s *Storage) User() storage.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		storage: s,
	}
	return s.userRepository
}
