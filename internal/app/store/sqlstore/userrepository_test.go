package sqlstore_test

import (
	"awesomeProject/internal/app/store"
	"awesomeProject/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDBHelper(t, false)
	defer teardown("users")

	s := sqlstore.NewStore(db)
	userGen := store.TestUserHelper(t)
	user := userGen()
	err := s.UserRepository().Create(user)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDBHelper(t, false)
	defer teardown("users")

	s := sqlstore.NewStore(db)

	email := "abc@gmail.com"

	user, err := s.UserRepository().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	userGen := store.TestUserHelper(t)
	user = userGen()
	err = s.UserRepository().Create(user)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	user, err = s.UserRepository().FindByEmail(email)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
}

func TestUserRepository_FindById(t *testing.T) {
	db, teardown := sqlstore.TestDBHelper(t, false)
	defer teardown("users")

	s := sqlstore.NewStore(db)
	userGen := store.TestUserHelper(t)
	user := userGen()
	err := s.UserRepository().Create(user)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	returnedUser, err := s.UserRepository().FindById(user.Id)
	assert.NoError(t, err)
	assert.Equal(t, returnedUser.Id, user.Id)
}
