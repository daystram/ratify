package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
)

func (m *module) AuthenticateUser(credentials datatransfers.UserLogin) (user models.User, sessionID string, err error) {
	if user, err = m.db.userOrmer.GetOneByUsername(credentials.Username); err != nil {
		return models.User{}, "", errors.New("incorrect credentials")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return models.User{}, "", errors.New("incorrect credentials")
	}
	sessionID = utils.GenerateRandomString(constants.SessionIDLength)
	sessionTokenKey := fmt.Sprintf(constants.RDKeySessionToken, sessionID)
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
	if result = m.rd.HGetAll(context.Background(), fmt.Sprintf(constants.RDKeySessionToken, sessionID)); result.Err() != nil {
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
