package datatransfers

type LogDetail struct {
	Scope  string      `json:"scope,omitempty"`
	Detail interface{} `json:"detail,omitempty"`
}

type LogInfo struct {
	Username        string `json:"preferred_username"`
	ApplicationName string `json:"application_name,omitempty"`
	ClientID        string `json:"client_id,omitempty"`
	Severity        string `json:"severity"`
	Description     string `json:"description"`
	CreatedAt       int64  `json:"created_at"`
}
