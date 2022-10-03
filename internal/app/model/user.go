package model

type Password struct {
	Original  string
	Encrypted string
}

type User struct {
	Id       int
	Email    string
	Password Password
}
