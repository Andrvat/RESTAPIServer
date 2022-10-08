package teststore

import (
	"awesomeProject/internal/app/model"
	"awesomeProject/internal/app/store"
)

type UserRepository struct {
	store     *Store
	usersById map[int]*model.User
}

func (r *UserRepository) Create(user *model.User) error {
	err := user.BeforeCreate()
	if err != nil {
		return err
	}

	user.Id = len(r.usersById) + 1
	r.usersById[user.Id] = user
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, user := range r.usersById {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

func (r *UserRepository) FindById(id int) (*model.User, error) {
	user, exist := r.usersById[id]
	if exist {
		return user, nil
	} else {
		return nil, store.ErrRecordNotFound
	}
}

func (r *UserRepository) GetAllUsers() ([]*model.User, error) {
	v := make([]*model.User, 0, len(r.usersById))

	for _, value := range r.usersById {
		v = append(v, value)
	}
	return v, nil
}
