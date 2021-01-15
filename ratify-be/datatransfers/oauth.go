package datatransfers

import (
	"github.com/daystram/ratify/ratify-be/constants"
)

type AuthorizationRequest struct {
	ResponseType string `json:"response_type" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	RedirectURI  string `json:"redirect_uri" binding:"-"`
	Scope        string `json:"scope" binding:"-"`
	State        string `json:"state" binding:"-"`
	UseSession   bool   `json:"use_session" binding:"-"`
	UserLogin
	PKCEAuthFields
}

func (authRequest *AuthorizationRequest) Flow() string {
	responseTypeCode := authRequest.ResponseType == constants.ResponseTypeCode
	hasCodeChallenge := authRequest.PKCEAuthFields.CodeChallenge != ""
	hasCodeChallengeMethod := authRequest.PKCEAuthFields.CodeChallengeMethod != ""
	switch {
	case responseTypeCode && !hasCodeChallenge && !hasCodeChallengeMethod:
		return constants.FlowAuthorizationCode
	case responseTypeCode && hasCodeChallenge && hasCodeChallengeMethod:
		return constants.FlowAuthorizationCodeWithPKCE
	default:
		return constants.FlowUnsupported
	}
}

type PKCEAuthFields struct {
	CodeChallenge       string `json:"code_challenge" binding:"-"`
	CodeChallengeMethod string `json:"code_challenge_method" binding:"-"`
}

type PKCETokenFields struct {
	CodeVerifier string `form:"code_verifier" binding:"-"`
}

type AuthorizationResponse struct {
	AuthorizationCode string `json:"code" url:"code"`
	State             string `json:"state" url:"state"`
}

type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required"`
	ClientSecret string `form:"client_secret" binding:"-"`
	Code         string `form:"code" binding:"required"`
	PKCETokenFields
}

func (tokenRequest *TokenRequest) Flow() string {
	grantTypeCode := tokenRequest.GrantType == constants.GrantTypeAuthorizationCode
	hasCodeVerifier := tokenRequest.PKCETokenFields.CodeVerifier != ""
	switch {
	case grantTypeCode && !hasCodeVerifier:
		return constants.FlowAuthorizationCode
	case grantTypeCode && hasCodeVerifier:
		return constants.FlowAuthorizationCodeWithPKCE
	default:
		return constants.FlowUnsupported
	}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type IntrospectRequest struct {
	Token        string `form:"token" binding:"required"`
	Hint         string `form:"token_type_hint" binding:"-"`
	ClientID     string `form:"client_id" binding:"required"`
	ClientSecret string `form:"client_secret" binding:"required"`
}

type TokenIntrospection struct {
	Active   bool   `json:"active"`
	ClientID string `json:"client_id,omitempty"`
	Subject  string `json:"sub,omitempty"`
	Scope    string `json:"scope,omitempty"`
}

type LogoutRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
	ClientID    string `json:"client_id" binding:"required"`
	Global      bool   `json:"global" binding:"-"`
}
