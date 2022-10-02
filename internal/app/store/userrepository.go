package store

import "awesomeProject/internal/app/model"

type UserRepository struct {
	store *Store
}

func (receiver UserRepository) Create(user *model.User) (*model.User, error) {
	err := receiver.store.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		user.Email,
		user.Password,
	).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (receiver UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := receiver.store.db.QueryRow(
		"SELECT id, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
