package constants

const (
	RDDelimiter      = "::"
	defaultDelimiter = RDDelimiter + "%s"

	RDKeyAuthorizationCode = GrantTypeAuthorizationCode + defaultDelimiter
	RDKeyCodeChallenge     = "code_challenge" + defaultDelimiter
	RDKeySessionToken      = "session_token" + defaultDelimiter
	RDKeyAccessToken       = "access_token" + defaultDelimiter
	RDKeyRefreshToken      = "refresh_token" + defaultDelimiter
)
