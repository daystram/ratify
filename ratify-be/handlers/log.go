package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
)

func (m *module) LogGetAllActivity(subject string) (logs []models.Log, err error) {
	if logs, err = m.db.logOrmer.GetAllByUserSubject(subject); err != nil {
		return nil, fmt.Errorf("cannot retrieve logs. %+v", err)
	}
	return
}

func (m *module) LogGetAllAdmin() (logs []models.Log, err error) {
	if logs, err = m.db.logOrmer.GetAllAdmin(); err != nil {
		return nil, fmt.Errorf("cannot retrieve logs. %+v", err)
	}
	return
}

func (m *module) LogInsertLogin(user models.User, application models.Application, success bool, detail datatransfers.LogDetail) {
	description, _ := json.Marshal(detail)
	m.logEntry(models.Log{
		UserSubject:         sql.NullString{String: user.Subject, Valid: true},
		ApplicationClientID: sql.NullString{String: application.ClientID, Valid: true},
		Type:                constants.LogTypeLogin,
		Severity:            map[bool]string{true: constants.LogSeverityInfo, false: constants.LogSeverityWarn}[success],
		Description:         string(description),
	})
}

func (m *module) LogInsertUser(user models.User, success bool, detail datatransfers.LogDetail) {
	description, _ := json.Marshal(detail)
	m.logEntry(models.Log{
		UserSubject: sql.NullString{String: user.Subject, Valid: true},
		Type:        constants.LogTypeUser,
		Severity:    map[bool]string{true: constants.LogSeverityInfo, false: constants.LogSeverityWarn}[success],
		Description: string(description),
	})
}

func (m *module) LogInsertApplication(user models.User, application models.Application, action bool, detail datatransfers.LogDetail) {
	description, _ := json.Marshal(detail)
	m.logEntry(models.Log{
		UserSubject:         sql.NullString{String: user.Subject, Valid: true},
		ApplicationClientID: sql.NullString{String: application.ClientID, Valid: !(detail.Scope == constants.LogScopeApplicationCreate && !action)},
		Type:                constants.LogTypeApplication,
		Severity:            map[bool]string{true: constants.LogSeverityInfo, false: constants.LogSeverityWarn}[action],
		Description:         string(description),
	})
}

func (m *module) logEntry(entry models.Log) {
	if err := m.db.logOrmer.InsertLog(entry); err != nil {
		log.Printf("failed adding log entry. %+v", err)
	}
}
