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

type TokenRequest struct {
	GrantType    string `json:"grant_type" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	Code         string `json:"code" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}
