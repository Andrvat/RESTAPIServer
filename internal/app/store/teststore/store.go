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
			store:     s,
			usersById: make(map[int]*model.User),
		}
	}
	return s.userRepository
}
