package constants

import (
	"time"
)

const (
	ResponseTypeCode = "code"

	AuthorizationCodeLength = 20
	AuthorizationCodeExpiry = time.Second * 30
)
