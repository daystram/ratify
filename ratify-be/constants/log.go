package constants

const (
	LogTypeLogin       = "LOGN"
	LogTypeUser        = "USER"
	LogTypeApplication = "APPN"

	LogSeverityFatal = "F"
	LogSeverityError = "E"
	LogSeverityWarn  = "W"
	LogSeverityInfo  = "I"

	LogScopeOAuthAuthorize = "oauth::authorize"
	LogScopeUserProfile    = "user::profile"
	LogScopeUserPassword   = "user::password"
	LogScopeUserMFA        = "user::mfa"
)
