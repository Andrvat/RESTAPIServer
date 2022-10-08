package store

import "awesomeProject/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindById(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
