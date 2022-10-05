package store

type Store interface {
	UserRepository() UserRepository
}
