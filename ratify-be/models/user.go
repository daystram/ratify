package models

import (
	"time"

	"github.com/daystram/ratify/ratify-be/constants"
	"gorm.io/gorm"
)

type userOrm struct {
	db *gorm.DB
}

// CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
type User struct {
	Subject       string `gorm:"column:sub;primaryKey;type:uuid;default:uuid_generate_v4()" json:"-"`
	Superuser     bool   `gorm:"column:superuser;default:false" json:"-"`
	GivenName     string `gorm:"column:given_name;type:varchar(20);not null" json:"-"`
	FamilyName    string `gorm:"column:family_name;type:varchar(12);not null" json:"-"`
	Username      string `gorm:"column:preferred_username;uniqueIndex;type:varchar(12);not null" json:"-"`
	Email         string `gorm:"column:email;uniqueIndex;type:varchar(50);not null" json:"-"`
	EmailVerified bool   `gorm:"column:email_verified;default:false" json:"-"`
	Password      string `gorm:"column:password;type:varchar(100);not null" json:"-"`
	TOTPSecret    string `gorm:"column:totp_secret;type:varchar(32);default:'-'" json:"-"`
	Metadata      string `gorm:"column:metadata;type:text" json:"-"`
	CreatedAt     int64  `gorm:"column:created_at;autoCreateTime" json:"-"`
	UpdatedAt     int64  `gorm:"column:updated_at;autoUpdateTime" json:"-"`
	LastLogin     int64  `gorm:"column:last_login;default:0" json:"-"`
	LoginCount    int    `gorm:"column:login_count;default:0" json:"-"`
}

type UserOrmer interface {
	GetOneBySubject(subject string) (user User, err error)
	GetOneByUsername(username string) (user User, err error)
	GetOneByEmail(email string) (user User, err error)
	InsertUser(user User) (subject string, err error)
	UpdateUser(user User) (err error)
	IncrementLoginCount(user User) (err error)
}

func NewUserOrmer(db *gorm.DB) UserOrmer {
	_ = db.AutoMigrate(&User{}) // builds table when enabled
	return &userOrm{db}
}

func (o *userOrm) GetOneBySubject(subject string) (user User, err error) {
	result := o.db.Model(&User{}).Where("sub	 = ?", subject).First(&user)
	return user, result.Error
}

func (o *userOrm) GetOneByUsername(username string) (user User, err error) {
	result := o.db.Model(&User{}).Where("preferred_username = ?", username).First(&user)
	return user, result.Error
}

func (o *userOrm) GetOneByEmail(email string) (user User, err error) {
	result := o.db.Model(&User{}).Where("email = ?", email).First(&user)
	return user, result.Error
}

func (o *userOrm) InsertUser(user User) (subject string, err error) {
	result := o.db.Model(&User{}).Create(&user)
	return user.Subject, result.Error
}

func (o *userOrm) UpdateUser(user User) (err error) {
	// By default, only non-empty fields are updated. See https://gorm.io/docs/update.html#Updates-multiple-columns
	result := o.db.Model(&User{}).Where("sub = ?", user.Subject).Updates(&user)
	return result.Error
}

func (o *userOrm) IncrementLoginCount(user User) (err error) {
	result := o.db.Model(&User{}).Where("sub = ?", user.Subject).
		Updates(map[string]interface{}{
			"last_login":  gorm.Expr("?", time.Now().Unix()),
			"login_count": gorm.Expr("login_count + 1"),
		})
	return result.Error
}

func (user User) EnabledTOTP() (enabled bool) {
	return user.TOTPSecret != constants.TOTPDisabledFlag
}
