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
		"ua_mobile":  userAgent.Mobile,
	}
	if err = m.rd.HSet(context.Background(), sessionIDKey, sessionIDValue).Err(); err != nil {
		return "", fmt.Errorf("failed storing session_id. %v", err)
	}
	if err = m.rd.Expire(context.Background(), sessionIDKey, constants.SessionIDExpiry).Err(); err != nil {
		m.rd.Del(context.Background(), sessionIDKey)
		return "", fmt.Errorf("failed setting session_id expiry. %v", err)
	}
	// initilize session_child (to allow setting expiry)
	sessionChildKey := fmt.Sprintf(constants.RDTemSessionChild, sessionID)
	if err = m.rd.ZAddNX(context.Background(), sessionChildKey, &redis.Z{
		Score:  -1,
		Member: "$$", // placeholder sentinel value
	}).Err(); err != nil {
		m.rd.Del(context.Background(), sessionIDKey)
		return "", fmt.Errorf("failed initializing session_child. %v", err)
	}
	if err = m.rd.Expire(context.Background(), sessionChildKey, constants.SessionIDExpiry).Err(); err != nil {
		m.rd.Del(context.Background(), sessionIDKey)
		m.rd.Del(context.Background(), sessionChildKey)
		return "", fmt.Errorf("failed setting session_child expiry. %v", err)
	}
	// collect session_id to list
	sessionListKey := fmt.Sprintf(constants.RDTemSessionList, subject)
	if err = m.rd.ZAddNX(context.Background(), sessionListKey, &redis.Z{
		Score:  float64(now + int64(constants.AccessTokenExpiry.Seconds())),
		Member: sessionID,
	}).Err(); err != nil {
		m.rd.Del(context.Background(), sessionIDKey)
		m.rd.Del(context.Background(), sessionChildKey)
		return "", fmt.Errorf("failed enlisting session_id. %v", err)
	}
	// prune dead sessions
	m.rd.ZRemRangeByScore(context.Background(), sessionListKey, "0", fmt.Sprintf("%d", now))
	return
}

func (m *module) SessionAddChild(sessionID, accessToken string) (err error) {
	// collect child access_token to list
	now := time.Now().Unix()
	sessionChildKey := fmt.Sprintf(constants.RDTemSessionChild, sessionID)
	if err = m.rd.ZAddNX(context.Background(), sessionChildKey, &redis.Z{
		Score:  float64(now + int64(constants.AccessTokenExpiry.Seconds())),
		Member: accessToken,
	}).Err(); err != nil {
		return fmt.Errorf("failed enlisting access_token. %v", err)
	}
	// prune dead children
	m.rd.ZRemRangeByScore(context.Background(), sessionChildKey, "0", fmt.Sprintf("%d", now))
	return
}

func (m *module) sessionDeleteChild(sessionID, accessToken string) (err error) {
	return m.rd.ZRem(context.Background(), fmt.Sprintf(constants.RDTemSessionChild, sessionID), accessToken).Err()
}

func (m *module) SessionInfo(sessionID string) (session datatransfers.SessionInfo, err error) {
	// retrieve session
	var result *redis.StringStringMapCmd
	if result = m.rd.HGetAll(context.Background(), fmt.Sprintf(constants.RDTemSessionID, sessionID)); result.Err() != nil {
		return datatransfers.SessionInfo{}, fmt.Errorf("failed retrieving session. %v", result.Err())
	}
	// build
	session.SessionID = sessionID
	session.Subject = result.Val()["subject"]
	session.UserAgent = datatransfers.UserAgent{
		IP:      result.Val()["ua_ip"],
		Browser: result.Val()["ua_browser"],
		OS:      result.Val()["ua_os"],
	}
	session.UserAgent.Mobile, _ = strconv.ParseBool(result.Val()["ua_mobile"])
	session.IssuedAt, _ = strconv.ParseInt(result.Val()["issued_at"], 10, 64)
	return session, nil
}

func (m *module) SessionRevoke(sessionID string) (err error) {
	// get session detail
	var session datatransfers.SessionInfo
	if session, err = m.SessionInfo(sessionID); err != nil {
		return fmt.Errorf("failed checking session_id. %v", err)
	}
	// remove session_id from session_list
	m.rd.ZRem(context.Background(), fmt.Sprintf(constants.RDTemSessionList, session.Subject), sessionID)
	// revoke session_id
	m.rd.Del(context.Background(), fmt.Sprintf(constants.RDTemSessionID, sessionID))
	// revoke associated access_token listed in session_child
	sessionChildKey := fmt.Sprintf(constants.RDTemSessionChild, sessionID)
	var result *redis.StringSliceCmd
	if result = m.rd.ZRangeByScore(context.Background(), sessionChildKey, &redis.ZRangeBy{
		Min: "0",
		Max: "+inf",
	}); result.Err() != nil {
		return fmt.Errorf("failed retrieving active sessions. %v", err)
	}
	for _, accessToken := range result.Val() {
		m.OAuthRevokeAccessToken(accessToken)
	}
	m.rd.Del(context.Background(), sessionChildKey)
	return
}

func (m *module) SessionGetAllActive(subject string) (activeSessions []*datatransfers.SessionInfo, err error) {
	// retrieve session_list
	var result *redis.StringSliceCmd
	if result = m.rd.ZRangeByScore(context.Background(), fmt.Sprintf(constants.RDTemSessionList, subject), &redis.ZRangeBy{
		Min: "0",
		Max: "+inf",
	}); result.Err() != nil {
		return nil, fmt.Errorf("failed retrieving active sessions. %v", err)
	}
	// build activeSessions
	activeSessions = make([]*datatransfers.SessionInfo, 0)
	for _, sessionID := range result.Val() {
		var session datatransfers.SessionInfo
		if session, err = m.SessionInfo(sessionID); err == nil {
			activeSessions = append(activeSessions, &session)
		}
	}
	return
}
