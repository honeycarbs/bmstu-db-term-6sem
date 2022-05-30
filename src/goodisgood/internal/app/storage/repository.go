package storage

import "goodisgood/internal/app/model"

type AccountRepository interface {
	Create(*model.Account) error
	FindByEmail(string) (*model.Account, error)
	Find(string) (*model.Account, error)
	GetAll() ([]model.Account, error)
}

type UserRepository interface {
	Create(string, *model.User) error
}

type LocationRepository interface {
	Create(*model.Location) error
	Assign(string, *model.Location) error
	Get() ([]model.Location, error)
}

type EducationPlaceRepository interface {
	Assign(string, *model.EducationPlace, *model.EducationProgram) error
	Get() ([]model.EducationPlace, error)
}

type PollRepository interface {
	Submit(string, *model.Answer) error
	GetUserResult(string) (*model.Poll, error)
	GetWordsList() ([]string, error)
	GetPollStats() ([]model.Stats, error)
}
