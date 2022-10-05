package teststore

import (
	"awesomeProject/internal/app/model"
	"awesomeProject/internal/app/store"
)

type Store struct {
	userRepository *UserRepository
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) UserRepository() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
			users: make(map[string]*model.User),
		}
	}
	return s.userRepository
}
