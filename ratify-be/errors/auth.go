package errors

import "errors"

var (
	ErrAuthIncorrectCredentials = errors.New("incorrect_credentials")
	ErrAuthEmailNotVerified     = errors.New("email_unverified")
)
