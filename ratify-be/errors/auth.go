package errors

import "errors"

var (
	ErrAuthIncorrectCredentials = errors.New("incorrect_credentials")
	ErrAuthMissingOTP           = errors.New("missing_otp")
	ErrAuthEmailNotVerified     = errors.New("email_unverified")
)
