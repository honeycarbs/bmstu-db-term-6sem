package neo4jstorage

import (
	"goodisgood/internal/app/storage"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Storage struct {
	db                       neo4j.Driver
	accountRepository        *AccountRepository
	userRepository           *UserRepository
	locationRepository       *LocationRepository
	educationPlaceRepository *EducationPlaceRepository
	pollRepository           *PollRepository
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

func (s *Storage) Location() storage.LocationRepository {
	if s.locationRepository != nil {
		return s.locationRepository
	}

	s.locationRepository = &LocationRepository{
		storage: s,
	}
	return s.locationRepository
}

func (s *Storage) EducationPlace() storage.EducationPlaceRepository {
	if s.educationPlaceRepository != nil {
		return s.educationPlaceRepository
	}
	s.educationPlaceRepository = &EducationPlaceRepository{
		storage: s,
	}
	return s.educationPlaceRepository
}

func (s *Storage) Poll() storage.PollRepository {
	if s.pollRepository != nil {
		return s.pollRepository
	}

	s.pollRepository = &PollRepository{
		storage: s,
	}
	return s.pollRepository
}
