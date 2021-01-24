package datatransfers

type DashboardInfo struct {
	SignInCount   int   `json:"signin_count"`
	LastSignIn    int64 `json:"last_signin"`
	SessionCount  int   `json:"session_count"`
	MFAEnabled    bool  `json:"mfa_enabled"`
	RecentFailure bool  `json:"recent_failure"`
}
