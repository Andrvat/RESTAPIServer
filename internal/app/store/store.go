package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	config         *Config
	db             *sql.DB
	userRepository *UserRepository
}

func NewStore(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (receiver *Store) Open() error {
	db, err := sql.Open(receiver.config.DatabaseDriverName, receiver.config.DatabaseUrl)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	receiver.db = db
	return nil
}

func (receiver Store) Close() error {
	err := receiver.db.Close()
	return err
}

func (receiver Store) UserRepository() *UserRepository {
	if receiver.userRepository == nil {
		receiver.userRepository = &UserRepository{
			store: &receiver,
		}
	}
	return receiver.userRepository
}
