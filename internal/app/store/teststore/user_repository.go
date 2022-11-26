package teststore

import (
	"awesomeProject/internal/app/model"
	"awesomeProject/internal/app/store"
	"errors"
)

type UserRepository struct {
	store     *Store
	usersById map[int]*model.User
}

func (r *UserRepository) Create(user *model.User) error {
	err := user.BeforeCreateOrUpdate()
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

func (r *UserRepository) AllUsers() ([]*model.User, error) {
	v := make([]*model.User, 0, len(r.usersById))

	for _, value := range r.usersById {
		v = append(v, value)
	}
	return v, nil
}

func (r *UserRepository) Update(user *model.User) error {
	_, exist := r.usersById[user.Id]
	if exist {
		r.usersById[user.Id] = user
		return nil
	} else {
		return errors.New("cannot update non-existing user")
	}
}

func (r *UserRepository) Delete(user *model.User) error {
	delete(r.usersById, user.Id)
	return nil
}
