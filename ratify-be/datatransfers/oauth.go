package datatransfers

import (
	"github.com/daystram/ratify/ratify-be/constants"
)

type AuthorizationRequest struct {
	ResponseType string `json:"response_type" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	RedirectURI  string `json:"redirect_uri" binding:"-"`
	State        string `json:"state" binding:"-"`
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
	CodeVerifier string `json:"code_verifier" binding:"-"`
}

type AuthorizationResponse struct {
	AuthorizationCode string `json:"code"`
	State             string `json:"state"`
}

type TokenRequest struct {
	GrantType    string `json:"grant_type" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"-"`
	Code         string `json:"code" binding:"required"`
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
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}
