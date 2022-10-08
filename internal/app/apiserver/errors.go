package apiserver

import "errors"

var (
	errIncorrectEmailOrPassword = errors.New("incorrect user email or password")
	errNotAuthenticated         = errors.New("user is not authenticated")
	errNonEmptyBodyRequired     = errors.New("server expected a non empty input body, but got null")
)
