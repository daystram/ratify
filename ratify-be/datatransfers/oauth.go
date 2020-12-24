package datatransfers

type AuthorizationRequest struct {
	ResponseType string `json:"response_type" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	RedirectURI  string `json:"redirect_uri" binding:"required"`
	State        string `json:"state" binding:"-"`
	UserLogin
}

type AuthorizationResponse struct {
	AuthorizationCode string `json:"code"`
	State             string `json:"state"`
}
