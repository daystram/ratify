package handlers

import (
	"errors"
	"time"

	"github.com/xlzd/gotp"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/models"
)

func (m *module) EnableTOTP(user models.User) (uri string, err error) {
	if user.EnabledTOTP() {
		return "", errors.New("totp already enabled")
	}
	secret := gotp.RandomSecret(constants.TOTPSecretLength)
	if err = m.db.userOrmer.UpdateUser(models.User{
		Subject:    user.Subject,
		TOTPSecret: secret,
	}); err != nil {
		return "", errors.New("failed storing totp_secret")
	}
	return gotp.NewDefaultTOTP(secret).ProvisioningUri(user.Username, constants.TOTPIssuer), nil
}

func (m *module) DisableTOTP(user models.User) (err error) {
	if !user.EnabledTOTP() {
		return errors.New("totp not yet enabled")
	}
	if err = m.db.userOrmer.UpdateUser(models.User{
		Subject:    user.Subject,
		TOTPSecret: constants.TOTPDisabledFlag,
	}); err != nil {
		return errors.New("failed deleting totp_secret")
	}
	return
}

func (m *module) CheckTOTP(otp string, user models.User) (valid bool) {
	return user.EnabledTOTP() && gotp.NewDefaultTOTP(user.TOTPSecret).Verify(otp, int(time.Now().Unix()))
}
