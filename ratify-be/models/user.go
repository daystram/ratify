package models

import (
	"gorm.io/gorm"
)

type userOrm struct {
	db *gorm.DB
}

// CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
type User struct {
	Subject       string `gorm:"column:sub;primaryKey;type:uuid;default:uuid_generate_v4()" json:"-"`
	Superuser     bool   `gorm:"column:superuser" json:"-"`
	Name          string `gorm:"column:name;type:varchar(50)" json:"-"`
	Username      string `gorm:"column:preferred_username;uniqueIndex;type:varchar(12)" json:"-"`
	Email         string `gorm:"column:email;uniqueIndex;type:varchar(50)" json:"-"`
	EmailVerified bool   `gorm:"column:email_verified" json:"-"`
	Password      string `gorm:"column:password;type:varchar(100)" json:"-"`
	Metadata      string `gorm:"column:metadata;type:text" json:"-"`
	CreatedAt     int64  `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt     int64  `gorm:"column:updated_at;autoUpdateTime" json:"-"`
}

type UserOrmer interface {
	GetOneBySubject(subject string) (user User, err error)
	GetOneByUsername(username string) (user User, err error)
	InsertUser(user User) (subject string, err error)
	UpdateUser(user User) (err error)
}

func NewUserOrmer(db *gorm.DB) UserOrmer {
	_ = db.AutoMigrate(&User{}) // builds table when enabled
	return &userOrm{db}
}

func (o *userOrm) GetOneBySubject(subject string) (user User, err error) {
	user.Subject = subject
	result := o.db.Model(&User{}).First(&user)
	return user, result.Error
}

func (o *userOrm) GetOneByUsername(username string) (user User, err error) {
	user.Username = username
	result := o.db.Model(&User{}).First(&user)
	return user, result.Error
}

func (o *userOrm) InsertUser(user User) (subject string, err error) {
	result := o.db.Model(&User{}).Create(&user)
	return user.Subject, result.Error
}

func (o *userOrm) UpdateUser(user User) (err error) {
	// By default, only non-empty fields are updated. See https://gorm.io/docs/update.html#Updates-multiple-columns
	result := o.db.Model(&User{}).Model(&user).Updates(&user)
	return result.Error
}
