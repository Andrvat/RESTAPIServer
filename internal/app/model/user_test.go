package model_test

import (
	"awesomeProject/internal/app/model"
	"awesomeProject/internal/app/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_BeforeCreate(t *testing.T) {
	userGen := store.TestUserHelper(t)
	user := userGen()
	err := user.BeforeCreate()
	assert.NoError(t, err)
	assert.NotEmpty(t, user.Password.Encrypted)
}

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		key     string
		userGen func() *model.User
		isValid bool
	}{
		{
			key:     "valid",
			userGen: store.TestUserHelper(t),
			isValid: true,
		},
		{
			key: "valid with encrypted password",
			userGen: func() *model.User {
				user := store.TestUserHelper(t)()
				user.Password.Original = ""
				user.Password.Encrypted = "hashCode"
				return user
			},
			isValid: true,
		},
		{
			key: "emptyEmail",
			userGen: func() *model.User {
				user := store.TestUserHelper(t)()
				user.Email = ""
				return user
			},
			isValid: false,
		},
		{
			key: "invalidEmail",
			userGen: func() *model.User {
				user := store.TestUserHelper(t)()
				user.Email = "xxxxyyyzzz"
				return user
			},
			isValid: false,
		},
		{
			key: "emptyPassword",
			userGen: func() *model.User {
				user := store.TestUserHelper(t)()
				user.Password.Original = ""
				return user
			},
			isValid: false,
		},
		{
			key: "shortPassword",
			userGen: func() *model.User {
				user := store.TestUserHelper(t)()
				user.Password.Original = "xyz"
				return user
			},
			isValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.key, func(t *testing.T) {
			if testCase.isValid {
				assert.NoError(t, testCase.userGen().Validate())
			} else {
				assert.Error(t, testCase.userGen().Validate())
			}
		})
	}
}
