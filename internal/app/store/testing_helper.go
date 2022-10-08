package store

import (
	"awesomeProject/internal/app/model"
	"testing"
)

func TestUserHelper(t *testing.T, optionalUserMeta ...interface{}) func() *model.User {
	t.Helper()

	return func() *model.User {
		switch len(optionalUserMeta) {
		case 3:
			return &model.User{
				Id:    optionalUserMeta[0].(int),
				Email: optionalUserMeta[1].(string),
				Password: &model.Password{
					Original: optionalUserMeta[2].(string),
				},
			}
		default:
			return &model.User{
				Id:       1,
				Email:    "abc@gmail.com",
				Password: &model.Password{Original: "super1234pass"},
			}
		}
	}
}
