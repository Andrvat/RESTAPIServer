package apiserver

import "errors"

var (
	ErrIncorrectEmailOrPassword = errors.New("incorrect user email or password")
	ErrNotAuthenticated         = errors.New("user is not authenticated")
	ErrNonEmptyBodyRequired     = errors.New("server expected a non empty input body, but got null")
)
