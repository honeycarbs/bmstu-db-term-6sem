package storage

type Storage interface {
	Account() AccountRepository
	User() UserRepository
	Location() LocationRepository
	EducationPlace() EducationPlaceRepository
	Poll() PollRepository
}
