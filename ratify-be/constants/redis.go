package constants

const (
	RDDelimiter      = "::"
	defaultDelimiter = RDDelimiter + "%s"

	RDKeyAuthorizationCode = GrantTypeAuthorizationCode + defaultDelimiter
	RDKeyCodeChallenge     = "code_challenge" + defaultDelimiter
	RDKeyAccessToken       = "access_token" + defaultDelimiter
	RDKeyRefreshToken      = "refresh_token" + defaultDelimiter
)
