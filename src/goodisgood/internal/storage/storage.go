package storage

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Storage struct {
	config            *Config
	db                neo4j.Driver
	accountRepository *AccountRepository
}

func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Open() error {
	db, err := neo4j.NewDriver(s.config.DBuri,
		neo4j.BasicAuth(s.config.DBUsername, s.config.DBPassword, ""))
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) Account() *AccountRepository {
	if s.accountRepository != nil {
		return s.accountRepository
	}

	s.accountRepository = &AccountRepository{
		storage: s,
	}

	return s.accountRepository
}
