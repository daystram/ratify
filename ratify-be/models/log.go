package models

import (
	"database/sql"
	"gorm.io/gorm"

	"github.com/daystram/ratify/ratify-be/constants"
)

type logOrm struct {
	db *gorm.DB
}

// CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
type Log struct {
	ID                  string `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"-"`
	User                User   `gorm:"foreignKey:UserSubject;references:Subject;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserSubject         sql.NullString
	Application         Application `gorm:"foreignKey:ApplicationClientID;default:null;references:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ApplicationClientID sql.NullString
	Type                string `gorm:"column:type;type:char(4);index:idx_type;not null" json:"-"`
	Severity            string `gorm:"column:severity;type:char(1);not null" json:"-"`
	Description         string `gorm:"column:description;type:text" json:"-"`
	CreatedAt           int64  `gorm:"column:created_at;autoCreateTime" json:"-"`
}

type LogOrmer interface {
	GetAllByUserSubject(userSubject string) (logs []Log, err error)
	GetAllActivityByApplicationClientID(applicationClientID string) (logs []Log, err error)
	InsertLog(log Log) (err error)
}

func NewLogOrmer(db *gorm.DB) LogOrmer {
	_ = db.AutoMigrate(&Log{}) // builds table when enabled
	return &logOrm{db}
}

func (o *logOrm) GetAllByUserSubject(userSubject string) (logs []Log, err error) {
	result := o.db.Model(&Log{}).
		Where("user_subject = ?", userSubject).
		Preload("User").Preload("Application").
		Order("created_at DESC").
		Find(&logs)
	return logs, result.Error
}

func (o *logOrm) GetAllActivityByApplicationClientID(applicationClientID string) (logs []Log, err error) {
	result := o.db.Model(&Log{}).
		Where("application_client_id = ? AND type = ?", applicationClientID, constants.LogTypeApplication).
		Preload("User").Preload("Application").
		Order("created_at DESC").
		Find(&logs)
	return logs, result.Error
}

func (o *logOrm) InsertLog(log Log) (err error) {
	result := o.db.Model(&Log{}).Create(&log)
	return result.Error
}
