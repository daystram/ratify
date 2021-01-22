package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	errors2 "github.com/daystram/ratify/ratify-be/errors"
	"github.com/daystram/ratify/ratify-be/models"
)

func (m *module) AuthenticateUser(credentials datatransfers.UserLogin) (user models.User, err error) {
	// check username/email
	if user, err = m.db.userOrmer.GetOneByUsername(credentials.Username); err != nil {
		if user, err = m.db.userOrmer.GetOneByEmail(credentials.Username); err != nil {
			return models.User{}, errors2.ErrAuthIncorrectIdentifier
		}
	}
	// check password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return user, errors2.ErrAuthIncorrectCredentials
	}
	// check email verification
	if !user.EmailVerified {
		return user, errors2.ErrAuthEmailNotVerified
	}
	// check if TOTP
	if user.EnabledTOTP() {
		if credentials.OTP == "" {
			return user, errors2.ErrAuthMissingOTP
		}
		if !m.CheckTOTP(credentials.OTP, user) {
			return user, errors2.ErrAuthIncorrectCredentials
		}
	}
	return user, nil
}

func (m *module) RegisterUser(userSignup datatransfers.UserSignup) (userSubject string, err error) {
	var hashedPassword []byte
	if hashedPassword, err = bcrypt.GenerateFromPassword([]byte(userSignup.Password), bcrypt.DefaultCost); err != nil {
		return "", errors.New("failed hashing password")
	}
	if userSubject, err = m.db.userOrmer.InsertUser(models.User{
		GivenName:  userSignup.GivenName,
		FamilyName: userSignup.FamilyName,
		Username:   userSignup.Username,
		Email:      userSignup.Email,
		Password:   string(hashedPassword),
	}); err != nil {
		return "", fmt.Errorf("error inserting user. %v", err)
	}
	return
}

func (m *module) VerifyUser(token string) (err error) {
	var result *redis.StringCmd
	if result = m.rd.Get(context.Background(), fmt.Sprintf(constants.RDTemVerificationToken, token)); result.Err() != nil {
		return fmt.Errorf("invalid verification_token. %v", result.Err())
	}
	_ = m.rd.Del(context.Background(), fmt.Sprintf(constants.RDTemVerificationToken, token))
	var user models.User
	if user, err = m.db.userOrmer.GetOneBySubject(result.Val()); err != nil {
		return fmt.Errorf("failed retrieving user. %v", result.Err())
	}
	user.EmailVerified = true
	if err = m.db.userOrmer.UpdateUser(user); err != nil {
		return fmt.Errorf("failed activating user. %v", result.Err())
	}
	return
}
