package constants

import (
	"time"
)

const (
	FlowAuthorizationCode         = "flow:authorization_code"
	FlowAuthorizationCodeWithPKCE = "flow:authorization_code_pkce"
	FlowUnsupported               = "flow:unsupported"

	ResponseTypeCode  = "code"
	ResponseTypeToken = "token"

	SessionIDCookieKey = "__Secure-sessid"

	PKCEChallengeMethodS256  = "S256"
	PKCEChallengeMethodPlain = "plain"

	GrantTypeAuthorizationCode = "authorization_code"

	AuthorizationCodeLength = 20
	AuthorizationCodeExpiry = time.Second * 300

	AccessTokenLength = 64
	AccessTokenExpiry = time.Hour * 2

	RefreshTokenLength = 64
	RefreshTokenExpiry = time.Hour * 24 * 14

	SessionIDLength = 64
	SessionIDExpiry = time.Hour * 10
)
