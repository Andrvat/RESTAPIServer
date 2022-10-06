package apiserver

import (
	"awesomeProject/internal/app/store/sqlstore"
	"database/sql"
	"log"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDatabaseConn(config.DatabaseUrl, config.DatabaseDriverName)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)

	store := sqlstore.NewStore(db)
	server := NewServer(store)
	err = http.ListenAndServe(config.BindAddr, server)
	return err
}

func newDatabaseConn(url string, driverName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
