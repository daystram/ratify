package models

import (
	"gorm.io/gorm"
)

type applicationOrm struct {
	db *gorm.DB
}

// CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
type Application struct {
	ID           string `gorm:"column:id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"-"`
	Owner        User   `gorm:"foreignKey:OwnerSubject;references:Subject;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	OwnerSubject string
	ClientID     string `gorm:"column:client_id;uniqueIndex;type:char(32);not null" json:"-"`
	ClientSecret string `gorm:"column:client_secret;type:char(64);not null" json:"-"`
	Name         string `gorm:"column:name;type:varchar(50);not null" json:"-"`
	Description  string `gorm:"column:description;type:varchar(255)" json:"-"`
	LoginURL     string `gorm:"column:login_url;type:text" json:"-"`
	CallbackURL  string `gorm:"column:callback_url;type:text" json:"-"`
	LogoutURL    string `gorm:"column:logout_url;type:text" json:"-"`
	Metadata     string `gorm:"column:metadata;type:text" json:"-"`
	CreatedAt    int64  `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt    int64  `gorm:"column:updated_at;autoUpdateTime" json:"-"`
}

type ApplicationOrmer interface {
	GetOneByClientID(clientID string) (application Application, err error)
	InsertApplication(application Application) (clientID string, err error)
	UpdateApplication(application Application) (err error)
}

func NewApplicationOrmer(db *gorm.DB) ApplicationOrmer {
	_ = db.AutoMigrate(&Application{}) // builds table when enabled
	return &applicationOrm{db}
}

func (o *applicationOrm) GetOneByClientID(clientID string) (application Application, err error) {
	result := o.db.Model(&Application{}).Where("client_id = ?", clientID).First(&application)
	return application, result.Error
}

func (o *applicationOrm) InsertApplication(application Application) (clientID string, err error) {
	result := o.db.Model(&Application{}).Create(&application)
	return application.ClientID, result.Error
}

func (o *applicationOrm) UpdateApplication(application Application) (err error) {
	result := o.db.Model(&User{}).Where("client_id = ?", application.ClientID).Updates(&application)
	return result.Error
}
