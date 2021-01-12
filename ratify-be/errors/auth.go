package errors

import "errors"

var (
	ErrAuthIncorrectIdentifier  = errors.New("incorrect_username")
	ErrAuthIncorrectCredentials = errors.New("incorrect_credentials")
	ErrAuthMissingOTP           = errors.New("missing_otp")
	ErrAuthEmailNotVerified     = errors.New("email_unverified")
)
