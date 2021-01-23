package constants

const (
	RDDelimiter      = "::"
	defaultDelimiter = RDDelimiter + "%s"

	RDTemAuthorizationCode = GrantTypeAuthorizationCode + defaultDelimiter
	RDTemCodeChallenge     = RDKeyCodeChallenge + defaultDelimiter
	RDTemSessionID         = RDKeySessionID + defaultDelimiter
	RDTemSessionList       = RDKeySessionList + defaultDelimiter
	RDTemSessionChild      = RDKeySessionChild + defaultDelimiter
	RDTemAccessToken       = RDKeyAccessToken + defaultDelimiter
	RDTemRefreshToken      = RDKeyRefreshToken + defaultDelimiter
	RDTemVerificationToken = RDKeyVerificationToken + defaultDelimiter
	RDTemTOTPToken         = RDKeyTOTPToken + defaultDelimiter

	RDKeyCodeChallenge     = "code_challenge"
	RDKeySessionID         = "session_id"
	RDKeySessionList       = "session_list"
	RDKeySessionChild      = "session_child"
	RDKeyAccessToken       = "access_token"
	RDKeyRefreshToken      = "refresh_token"
	RDKeyVerificationToken = "refresh_token"
	RDKeyTOTPToken         = "totp_token"
)
