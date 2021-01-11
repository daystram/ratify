package constants

import "time"

const (
	TOTPSecretLength  = 32
	TOTPConfirmExpiry = time.Minute * 10
	TOTPIssuer        = "Ratify"
	TOTPDisabledFlag  = "-"
)
