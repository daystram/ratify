package datatransfers

type UserAgent struct {
	IP      string `json:"ip"`
	Browser string `json:"browser"`
	OS      string `json:"os"`
	Mobile  bool   `json:"mobile"`
}

type Session struct {
	SessionID string `json:"session_id"`
	Subject   string `json:"-"`
	UserAgent
	IssuedAt int64 `json:"issued_at"`
}
