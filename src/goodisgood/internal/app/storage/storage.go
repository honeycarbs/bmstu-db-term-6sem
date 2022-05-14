package storage

type Storage interface {
	Account() AccountRepository
	User() UserRepository
}
