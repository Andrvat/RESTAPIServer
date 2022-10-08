package store

import (
	"awesomeProject/internal/app/model"
	"testing"
)

func TestUserHelper(t *testing.T) func() *model.User {
	t.Helper()

	return func() *model.User {
		return &model.User{
			Id:       1,
			Email:    "abc@gmail.com",
			Password: &model.Password{Original: "super1234pass"},
		}
	}
}
