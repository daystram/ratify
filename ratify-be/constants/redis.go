package constants

const (
	RDDelimiter      = "::"
	defaultDelimiter = RDDelimiter + "%s"

	RDTemAuthorizationCode = GrantTypeAuthorizationCode + defaultDelimiter
	RDTemCodeChallenge     = RDKeyCodeChallenge + defaultDelimiter
	RDTemSessionToken      = RDKeySessionToken + defaultDelimiter
	RDTemAccessToken       = RDKeyAccessToken + defaultDelimiter
	RDTemRefreshToken      = RDKeyRefreshToken + defaultDelimiter

	RDKeyCodeChallenge = "code_challenge"
	RDKeySessionToken  = "session_token"
	RDKeyAccessToken   = "access_token"
	RDKeyRefreshToken  = "refresh_token"
)
