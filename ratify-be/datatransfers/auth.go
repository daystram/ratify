package datatransfers

import (
	"errors"
	"time"
)

type OpenIDClaims struct {
	Subject   string `json:"sub"`
	Issuer    string `json:"iss"`
	Audience  string `json:"aud"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`

	// scope: profile
	Superuser  *bool   `json:"is_superuser,omitempty"`
	Username   *string `json:"preferred_username,omitempty"`
	GivenName  *string `json:"given_name,omitempty"`
	FamilyName *string `json:"family_name,omitempty"`
	UpdatedAt  *int64  `json:"updated_at,omitempty"`
	CreatedAt  *int64  `json:"created_at,omitempty"`

	// scope: email
	Email         *string `json:"email,omitempty"`
	EmailVerified *bool   `json:"email_verified,omitempty"`
}

func (c OpenIDClaims) Valid() (err error) {
	now := time.Now()
	if now.After(time.Unix(c.ExpiresAt, 0)) {
		err = errors.New("token has expired")
	}
	if now.Before(time.Unix(c.IssuedAt, 0)) {
		err = errors.New("token used before issued")
	}
	return err
}

type JWTClaims struct {
	Subject     string `json:"sub"`
	IsSuperuser bool   `json:"sup"`
	ExpiresAt   int64  `json:"exp"`
	IssuedAt    int64  `json:"iat"`
}

func (c JWTClaims) Valid() (err error) {
	now := time.Now()
	if now.After(time.Unix(c.ExpiresAt, 0)) {
		err = errors.New("token has expired")
	}
	if now.Before(time.Unix(c.IssuedAt, 0)) {
		err = errors.New("token used before issued")
	}
	return err
}
