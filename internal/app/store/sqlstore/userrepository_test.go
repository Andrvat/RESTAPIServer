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

func TestUserRepository_GetAllUsers(t *testing.T) {
	db, teardown := sqlstore.TestDBHelper(t, false)
	defer teardown("users")

	s := sqlstore.NewStore(db)
	userGen := store.TestUserHelper(t, 1, "abcabcabc@mail.com", "1234567890")
	user := userGen()
	err := s.UserRepository().Create(user)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	userGen = store.TestUserHelper(t, 2, "abcabc@mail.com", "1234567890")
	user = userGen()
	err = s.UserRepository().Create(user)

	users, err := s.UserRepository().GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
}

func TestUserRepository_Update(t *testing.T) {
	db, teardown := sqlstore.TestDBHelper(t, false)
	defer teardown("users")

	s := sqlstore.NewStore(db)
	userGen := store.TestUserHelper(t)
	user := userGen()
	err := s.UserRepository().Create(user)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	newEmail := "abababa@mail.com"
	user.Email = newEmail
	err = s.UserRepository().Update(user)
	assert.NoError(t, err)
	assert.Equal(t, newEmail, user.Email)
}
