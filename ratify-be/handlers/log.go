package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"encoding/json"
	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
)

func (m *module) RetrieveActivityLogs(subject string) (logs []models.Log, err error) {
	if logs, err = m.db.logOrmer.GetAllByUserSubject(subject); err != nil {
		return nil, errors.New(fmt.Sprintf("cannot retrieve logs. %+v", err))
	}
	return
}

func (m *module) LogLogin(user models.User, application models.Application, success bool, detail datatransfers.LogDetail) {
	description, _ := json.Marshal(detail)
	m.logEntry(models.Log{
		UserSubject:         sql.NullString{String: user.Subject, Valid: true},
		ApplicationClientID: sql.NullString{String: application.ClientID, Valid: true},
		Type:                constants.LogTypeLogin,
		Severity:            map[bool]string{true: constants.LogSeverityInfo, false: constants.LogSeverityWarn}[success],
		Description:         string(description),
	})
}

func (m *module) LogUser(user models.User, success bool, detail datatransfers.LogDetail) {
	description, _ := json.Marshal(detail)
	m.logEntry(models.Log{
		UserSubject: sql.NullString{String: user.Subject, Valid: true},
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
