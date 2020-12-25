package constants

import (
	"time"
)

const (
	ResponseTypeCode = "code"

	GrantTypeAuthorizationCode = "authorization_code"

	AuthorizationCodeLength = 20
	AuthorizationCodeExpiry = time.Second * 300

	AccessTokenLength  = 64
	AccessTokenExpiry  = time.Hour * 10
	RefreshTokenLength = 64
	RefreshTokenExpiry = time.Hour * 24 * 14
)
