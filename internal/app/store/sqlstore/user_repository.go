package sqlstore

import (
	"awesomeProject/internal/app/model"
	"awesomeProject/internal/app/store"
	"database/sql"
	"log"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(user *model.User) error {
	err := user.BeforeCreateOrUpdate()
	if err != nil {
		return err
	}
	err = r.store.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		user.Email,
		user.Password.Encrypted,
	).Scan(&user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := model.NewEmptyUser()
	err := r.store.db.QueryRow(
		"SELECT id, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.Id, &user.Email, &user.Password.Encrypted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindById(id int) (*model.User, error) {
	user := model.NewEmptyUser()
	err := r.store.db.QueryRow(
		"SELECT id, email, password FROM users WHERE id = $1",
		id).Scan(&user.Id, &user.Email, &user.Password.Encrypted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetAllUsers() ([]*model.User, error) {
	rows, err := r.store.db.Query("SELECT id, email FROM users")
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Query didn't close correctly")
		}
	}(rows)

	if err != nil {
		return nil, store.ErrRecordNotFound
	}

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err = rows.Scan(&user.Id, &user.Email)
		if err != nil {
			return nil, store.ErrDatabaseInternal
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) Update(user *model.User) error {
	err := user.BeforeCreateOrUpdate()
	if err != nil {
		return err
	}
	_, err = r.store.db.Exec("UPDATE users SET email = $2, password = $3 WHERE id = $1", user.Id, user.Email, user.Password.Encrypted)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(user *model.User) error {
	_, err := r.store.db.Exec("DELETE FROM users WHERE id = $1", user.Id)
	if err != nil {
		return err
	}
	return nil
}
