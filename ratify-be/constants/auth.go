package constants

import (
	"time"
)

const (
	EmailVerificationTokenLength = 32
	EmailVerificationTokenExpiry = time.Minute * 30
)
