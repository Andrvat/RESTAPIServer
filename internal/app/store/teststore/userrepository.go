package teststore

import (
	"awesomeProject/internal/app/model"
	"awesomeProject/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (r *UserRepository) Create(user *model.User) error {
	err := user.BeforeCreate()
	if err != nil {
		return err
	}

	r.users[user.Email] = user
	user.Id = len(r.users)
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	user, exist := r.users[email]
	if exist {
		return user, nil
	} else {
		return nil, store.ErrRecordNotFound
	}
}
