package constants

const (
	defaultSeparator = "::%s"

	AuthorizationCodeKey = GrantTypeAuthorizationCode + defaultSeparator
	AccessTokenKey       = "access_token" + defaultSeparator
	RefreshTokenKey      = "refresh_token" + defaultSeparator
)
