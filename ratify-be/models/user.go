package models

import (
	"gorm.io/gorm"
)

type userOrm struct {
	db *gorm.DB
}

type User struct {
	ID            string `gorm:"column:sub;primaryKey" json:"-"`
	Superuser     bool   `gorm:"column:superuser" json:"-"`
	Name          string `gorm:"column:name" json:"-"`
	Username      string `gorm:"column:preferred_username;uniqueIndex" json:"-"`
	Email         string `gorm:"column:email;uniqueIndex" json:"-"`
	EmailVerified bool   `gorm:"column:email_verified" json:"-"`
	Password      string `gorm:"column:password" json:"-"`
	Metadata      string `gorm:"column:metadata" json:"-"`
	CreatedAt     int64  `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt     int64  `gorm:"column:updated_at;autoUpdateTime" json:"-"`
}

type UserOrmer interface {
	GetOneByID(id string) (user User, err error)
	GetOneByUsername(username string) (user User, err error)
	InsertUser(user User) (id string, err error)
	UpdateUser(user User) (err error)
}

func NewUserOrmer(db *gorm.DB) UserOrmer {
	//_ = db.AutoMigrate(&User{})		// builds table when enabled
	return &userOrm{db}
}

func (o *userOrm) GetOneByID(id string) (user User, err error) {
	user.ID = id
	result := o.db.Model(&User{}).First(&user)
	return user, result.Error
}

func (o *userOrm) GetOneByUsername(username string) (user User, err error) {
	user.Username = username
	result := o.db.Model(&User{}).First(&user)
	return user, result.Error
}

func (o *userOrm) InsertUser(user User) (id string, err error) {
	result := o.db.Model(&User{}).Create(&user)
	return user.ID, result.Error
}

func (o *userOrm) UpdateUser(user User) (err error) {
	// By default, only non-empty fields are updated. See https://gorm.io/docs/update.html#Updates-multiple-columns
	result := o.db.Model(&User{}).Model(&user).Updates(&user)
	return result.Error
}
