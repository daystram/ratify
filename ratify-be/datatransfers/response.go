package datatransfers

type APIResponse struct {
	Code  interface{} `json:"code,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}
