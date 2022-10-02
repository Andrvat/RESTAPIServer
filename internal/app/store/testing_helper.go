package store

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestStoreHelper(t *testing.T, readEnv bool) (*Store, func(...string)) {
	t.Helper()

	config := NewConfig()
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
	config.DatabaseUrl = databaseUrl
	config.DatabaseDriverName = databaseDriverName

	store := NewStore(config)
	err := store.Open()
	if err != nil {
		t.Fatal(err)
	}

	return store, func(tables ...string) {
		if len(tables) > 0 {
			_, err := store.db.Exec(fmt.Sprintf(
				"TRUNCATE %s CASCADE", strings.Join(tables, ", "),
			))
			if err != nil {
				t.Fatal(err)
			}
		}
		err := store.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
}
