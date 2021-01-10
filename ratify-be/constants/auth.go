package constants

import (
	"time"
)

const (
	AuthenticationTimeout = time.Hour * 24 * 2

	EmailVerificationTokenLength = 32
	EmailVerificationTokenExpiry = time.Minute * 30
)
