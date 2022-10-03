package store_test

import (
	"awesomeProject/internal/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := store.TestStoreHelper(t, false)
	defer teardown("users")

	userGen := store.TestUserHelper(t)
	user, err := s.UserRepository().Create(userGen())
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s, teardown := store.TestStoreHelper(t, false)
	defer teardown("users")

	email := "abc@gmail.com"

	user, err := s.UserRepository().FindByEmail(email)
	assert.Error(t, err)

	userGen := store.TestUserHelper(t)
	user, err = s.UserRepository().Create(userGen())
	assert.NoError(t, err)
	assert.NotNil(t, user)

	user, err = s.UserRepository().FindByEmail(email)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
}
