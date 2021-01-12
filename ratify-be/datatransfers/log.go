package datatransfers

type LogDetail struct {
	Scope  string      `json:"scope,omitempty"`
	Detail interface{} `json:"detail,omitempty"`
}
