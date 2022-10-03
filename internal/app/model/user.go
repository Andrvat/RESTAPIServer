package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Original  string
	Encrypted string
}

type User struct {
	Id       int
	Email    string
	Password Password
}

func (u *User) BeforeCreate() error {
	err := u.Validate()
	if err != nil {
		return err
	}
	err = u.encryptPassword()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) encryptPassword() error {
	if len(u.Password.Original) > 0 {
		encrypted, err := encrypt(u.Password.Original)
		if err != nil {
			return err
		}
		u.Password.Encrypted = encrypted
	}
	return nil
}

func encrypt(original string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(original), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}

func (u *User) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password),
	)
	return err
}

func (p Password) Validate() error {
	err := validation.ValidateStruct(&p,
		validation.Field(&p.Original,
			validation.By(func(condition bool) validation.RuleFunc {
				return func(value interface{}) error {
					if condition {
						return validation.Validate(value, validation.Required)
					}
					return nil
				}
			}(p.Encrypted == "")),
			validation.Length(8, 36)),
	)
	return err
}
