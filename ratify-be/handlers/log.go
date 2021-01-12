package handlers

import (
	"log"

	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
)

func (m *module) LogLogin(user models.User, application models.Application, success bool, detail datatransfers.LogDetail) {
	description, _ := json.Marshal(detail)
	m.logEntry(models.Log{
		UserSubject:         user.Subject,
		ApplicationClientID: application.ClientID,
		Type:                constants.LogTypeLogin,
		Severity:            map[bool]string{true: constants.LogSeverityInfo, false: constants.LogSeverityWarn}[success],
		Description:         string(description),
	})
}

func (m *module) LogUser(user models.User, success bool, detail datatransfers.LogDetail) {
	description, _ := json.Marshal(detail)
	m.logEntry(models.Log{
		UserSubject: user.Subject,
		Type:        constants.LogTypeUser,
		Severity:    map[bool]string{true: constants.LogSeverityInfo, false: constants.LogSeverityWarn}[success],
		Description: string(description),
	})
}

func (m *module) logEntry(entry models.Log) {
	if err := m.db.logOrmer.InsertLog(entry); err != nil {
		log.Printf("failed adding log entry. %+v", err)
	}
}

func (m *module) ParseUserAgent(c *gin.Context) (ip, browser, os string) {
	ua := user_agent.New(c.Request.UserAgent())
	browser, _ = ua.Browser()
	return c.ClientIP(), browser, ua.OS()
}
