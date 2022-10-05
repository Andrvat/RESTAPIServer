package sqlstore

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestDBHelper(t *testing.T, readEnv bool) (*sql.DB, func(...string)) {
	t.Helper()

	var databaseUrl string
	var databaseDriverName string
	if readEnv {
		databaseUrl = os.Getenv("DATABASE_URL")
		databaseDriverName = os.Getenv("DATABASE_DRIVER_NAME")
	}
	if databaseUrl == "" || !readEnv {
		databaseUrl = "host=localhost port=5432 user=andrvat password=1234 dbname=awesome_project_dev_test"
	}
	if databaseDriverName == "" || !readEnv {
		databaseDriverName = "postgres"
	}

	db, err := sql.Open(databaseDriverName, databaseUrl)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			_, err := db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
			if err != nil {
				t.Fatal(err)
			}
		}
		err := db.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
}
