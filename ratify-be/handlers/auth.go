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
	"github.com/daystram/ratify/ratify-be/utils"
)

func (m *module) AuthenticateUser(credentials datatransfers.UserLogin) (user models.User, sessionID string, err error) {
	if user, err = m.db.userOrmer.GetOneByUsername(credentials.Username); err != nil {
		return models.User{}, "", errors2.ErrAuthIncorrectCredentials
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return models.User{}, "", errors2.ErrAuthIncorrectCredentials
	}
	if !user.EmailVerified {
		return models.User{}, "", errors2.ErrAuthEmailNotVerified
	}
	sessionID = utils.GenerateRandomString(constants.SessionIDLength)
	sessionTokenKey := fmt.Sprintf(constants.RDTemSessionToken, sessionID)
	sessionTokenValue := map[string]interface{}{
		"subject": user.Subject,
	}
	if err = m.rd.HSet(context.Background(), sessionTokenKey, sessionTokenValue).Err(); err != nil {
		return models.User{}, "", errors.New(fmt.Sprintf("failed storing session_id. %v", err))
	}
	if err = m.rd.Expire(context.Background(), sessionTokenKey, constants.SessionIDExpiry).Err(); err != nil {
		m.rd.Del(context.Background(), sessionTokenKey)
		return models.User{}, "", errors.New(fmt.Sprintf("failed setting session_id expiry. %v", err))
	}
	return user, sessionID, nil
}

func (m *module) CheckSession(sessionID string) (user models.User, newSessionID string, err error) {
	var result *redis.StringStringMapCmd
	if result = m.rd.HGetAll(context.Background(), fmt.Sprintf(constants.RDTemSessionToken, sessionID)); result.Err() != nil {
		return models.User{}, "", errors.New(fmt.Sprintf("failed retrieving authorization_code. %v", result.Err()))
	}
	// TODO: rotate/refresh sessionID?
	var userSubject string
	userSubject = result.Val()["subject"]
	if user, err = m.db.userOrmer.GetOneBySubject(userSubject); err != nil {
		return models.User{}, "", errors.New(fmt.Sprintf("failed retrieving user. %v", result.Err()))
	}
	return user, sessionID, nil
}

func (m *module) ClearSession(sessionID string) (err error) {
	return m.rd.Del(context.Background(), fmt.Sprintf(constants.RDTemSessionToken, sessionID)).Err()
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
		return "", errors.New(fmt.Sprintf("error inserting user. %v", err))
	}
	return
}

func (m *module) VerifyUser(token string) (err error) {
	var result *redis.StringCmd
	if result = m.rd.Get(context.Background(), fmt.Sprintf(constants.RDTemVerificationToken, token)); result.Err() != nil {
		return errors.New(fmt.Sprintf("invalid verification_token. %v", result.Err()))
	}
	_ = m.rd.Del(context.Background(), fmt.Sprintf(constants.RDTemVerificationToken, token))
	var user models.User
	if user, err = m.db.userOrmer.GetOneBySubject(result.Val()); err != nil {
		return errors.New(fmt.Sprintf("failed retrieving user. %v", result.Err()))
	}
	user.EmailVerified = true
	if err = m.db.userOrmer.UpdateUser(user); err != nil {
		return errors.New(fmt.Sprintf("failed activating user. %v", result.Err()))
	}
	return
}
