package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/utils"
)

func (m *module) SessionInitialize(subject string, userAgent datatransfers.UserAgent) (sessionID string, err error) {
	// generate session_id
	now := time.Now().Unix()
	sessionID = utils.GenerateRandomString(constants.SessionIDLength)
	sessionIDKey := fmt.Sprintf(constants.RDTemSessionID, sessionID)
	sessionIDValue := map[string]interface{}{
		"subject":    subject,
		"issued_at":  now,
		"ua_ip":      userAgent.IP,
		"ua_browser": userAgent.Browser,
		"ua_os":      userAgent.OS,
	}
	if err = m.rd.HSet(context.Background(), sessionIDKey, sessionIDValue).Err(); err != nil {
		return "", fmt.Errorf("failed storing session_id. %v", err)
	}
	if err = m.rd.Expire(context.Background(), sessionIDKey, constants.SessionIDExpiry).Err(); err != nil {
		m.rd.Del(context.Background(), sessionIDKey)
		return "", fmt.Errorf("failed setting session_id expiry. %v", err)
	}
	// collect session_id to list
	sessionListKey := fmt.Sprintf(constants.RDTemSessionList, subject)
	if err = m.rd.ZAdd(context.Background(), sessionListKey, &redis.Z{
		Score:  float64(now + int64(constants.AccessTokenExpiry.Seconds())),
		Member: sessionID,
	}).Err(); err != nil {
		m.rd.Del(context.Background(), sessionIDKey)
		return "", fmt.Errorf("failed listing session_id. %v", err)
	}
	m.rd.ZRemRangeByScore(context.Background(), sessionListKey, "0", fmt.Sprintf("%d", now))
	return
}

func (m *module) SessionAddChild(sessionID, accessToken string) (err error) {
	// collect child access_token to list
	now := time.Now().Unix()
	sessionChildKey := fmt.Sprintf(constants.RDTemSessionChild, sessionID)
	if err = m.rd.ZAdd(context.Background(), sessionChildKey, &redis.Z{
		Score:  float64(now + int64(constants.AccessTokenExpiry.Seconds())),
		Member: accessToken,
	}).Err(); err != nil {
		return fmt.Errorf("failed listing access_token. %v", err)
	}
	m.rd.ZRemRangeByScore(context.Background(), sessionChildKey, "0", fmt.Sprintf("%d", now))
	return
}

func (m *module) SessionInfo(sessionID string) (session datatransfers.Session, err error) {
	// retrieve session
	var result *redis.StringStringMapCmd
	if result = m.rd.HGetAll(context.Background(), fmt.Sprintf(constants.RDTemSessionID, sessionID)); result.Err() != nil {
		return datatransfers.Session{}, fmt.Errorf("failed retrieving session. %v", result.Err())
	}
	// build
	session.SessionID = sessionID
	session.Subject = result.Val()["subject"]
	session.UserAgent = datatransfers.UserAgent{
		IP:      result.Val()["ua_ip"],
		Browser: result.Val()["ua_browser"],
		OS:      result.Val()["ua_os"],
	}
	session.IssuedAt, _ = strconv.ParseInt(result.Val()["issued_at"], 10, 64)
	return session, nil
}

func (m *module) SessionClear(sessionID string) (err error) {
	// get session detail
	var session datatransfers.Session
	if session, err = m.SessionInfo(sessionID); err != nil {
		return fmt.Errorf("failed checking session_id. %v", err)
	}
	// remove session_id from list
	sessionListKey := fmt.Sprintf(constants.RDTemSessionList, session.Subject)
	if err = m.rd.ZRem(context.Background(), sessionListKey, sessionID).Err(); err != nil {
		return fmt.Errorf("failed unlisting session_id. %v", err)
	}
	return m.rd.Del(context.Background(), fmt.Sprintf(constants.RDTemSessionID, sessionID)).Err()
}

func (m *module) GetActiveSessions(subject string) (activeSessions []datatransfers.Session, err error) {
	// retrieve session_list
	var result *redis.StringSliceCmd
	if result = m.rd.ZRangeByScore(context.Background(), fmt.Sprintf(constants.RDTemSessionList, subject), &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", time.Now().Unix()),
		Max: "+inf",
	}); result.Err() != nil {
		return nil, fmt.Errorf("failed retrieving active sessions. %v", err)
	}
	// build activeSessions
	activeSessions = make([]datatransfers.Session, 0)
	for _, sessionID := range result.Val() {
		var session datatransfers.Session
		if session, err = m.SessionInfo(sessionID); err == nil {
			activeSessions = append(activeSessions, session)
		}
	}
	return
}
