package teststore_test

import (
	"awesomeProject/internal/app/store"
	"awesomeProject/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.NewStore()
	userGen := store.TestUserHelper(t)
	user := userGen()
	err := s.UserRepository().Create(user)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.NewStore()

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