package handlers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
)

func (m *module) SessionInitialize(subject string) (sessionID string, err error) {
	// generate session_id
	sessionID = utils.GenerateRandomString(constants.SessionIDLength)
	sessionIDKey := fmt.Sprintf(constants.RDTemSessionID, sessionID)
	sessionIDValue := map[string]interface{}{
		"subject": subject,
	}
	if err = m.rd.HSet(context.Background(), sessionIDKey, sessionIDValue).Err(); err != nil {
		return "", fmt.Errorf("failed storing session_id. %v", err)
	}
	if err = m.rd.Expire(context.Background(), sessionIDKey, constants.SessionIDExpiry).Err(); err != nil {
		m.rd.Del(context.Background(), sessionIDKey)
		return "", fmt.Errorf("failed setting session_id expiry. %v", err)
	}
	return
}


func (m *module) SessionCheck(sessionID string) (user models.User, newSessionID string, err error) {
	var result *redis.StringStringMapCmd
	if result = m.rd.HGetAll(context.Background(), fmt.Sprintf(constants.RDTemSessionID, sessionID)); result.Err() != nil {
		return models.User{}, "", fmt.Errorf("failed retrieving session. %v", result.Err())
	}
	// TODO: rotate/refresh sessionID?
	var userSubject string
	userSubject = result.Val()["subject"]
	if user, err = m.db.userOrmer.GetOneBySubject(userSubject); err != nil {
		return models.User{}, "", fmt.Errorf("failed retrieving user. %v", result.Err())
	}
	return user, sessionID, nil
}

func (m *module) SessionClearCurrent(sessionID string) (err error) {
	return m.rd.Del(context.Background(), fmt.Sprintf(constants.RDTemSessionID, sessionID)).Err()
}

