package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xlzd/gotp"

	"github.com/daystram/ratify/ratify-be/constants"
	errors2 "github.com/daystram/ratify/ratify-be/errors"
	"github.com/daystram/ratify/ratify-be/models"
)

func (m *module) EnableTOTP(user models.User) (uri string, err error) {
	if user.EnabledTOTP() {
		return "", errors.New("totp already enabled")
	}
	secret := gotp.RandomSecret(constants.TOTPSecretLength)
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf(constants.RDTemTOTPToken, user.Subject), secret, constants.TOTPConfirmExpiry).Err(); err != nil {
		return "", errors.New(fmt.Sprintf("failed storing totp_token. %v", err))
	}
	return gotp.NewDefaultTOTP(secret).ProvisioningUri(user.Username, constants.TOTPIssuer), nil
}

func (m *module) ConfirmTOTP(otp string, user models.User) (err error) {
	var result *redis.StringCmd
	if result = m.rd.Get(context.Background(), fmt.Sprintf(constants.RDTemTOTPToken, user.Subject)); result.Err() != nil {
		return errors2.ErrAuthIncorrectCredentials
	}
	secret := result.Val()
	if !gotp.NewDefaultTOTP(secret).Verify(otp, int(time.Now().Unix())) {
		return errors2.ErrAuthIncorrectCredentials
	}
	_ = m.rd.Del(context.Background(), fmt.Sprintf(constants.RDTemTOTPToken, user.Subject))
	if err = m.db.userOrmer.UpdateUser(models.User{
		Subject:    user.Subject,
		TOTPSecret: secret,
	}); err != nil {
		return errors.New("failed storing totp_secret")
	}
	return
}

func (m *module) DisableTOTP(user models.User) (err error) {
	if !user.EnabledTOTP() {
		return errors.New("totp not enabled")
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
