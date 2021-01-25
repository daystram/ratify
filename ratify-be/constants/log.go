package constants

const (
	LogTypeLogin       = "LOGN"
	LogTypeUserAdmin   = "USAD"
	LogTypeUser        = "USER"
	LogTypeApplication = "APPN"

	LogSeverityFatal = "F"
	LogSeverityError = "E"
	LogSeverityWarn  = "W"
	LogSeverityInfo  = "I"

	LogScopeOAuthAuthorize    = "oauth::authorize"
	LogScopeUserProfile       = "user::profile"
	LogScopeUserPassword      = "user::password"
	LogScopeUserSuperuser     = "user::superuser"
	LogScopeUserSession       = "user::session"
	LogScopeUserMFA           = "user::mfa"
	LogScopeApplicationDetail = "application::detail"
	LogScopeApplicationCreate = "application::create"
	LogScopeApplicationSecret = "application::secret"
)
