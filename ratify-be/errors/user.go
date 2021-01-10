package errors

import "errors"

var (
	ErrUserEmailExists    = errors.New("email_exists")
	ErrUserUsernameExists = errors.New("username_exists")
)
