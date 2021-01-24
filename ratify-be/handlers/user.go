package handlers

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/daystram/ratify/ratify-be/datatransfers"
	errors2 "github.com/daystram/ratify/ratify-be/errors"
	"github.com/daystram/ratify/ratify-be/models"
)

func (m *module) UserGetOneBySubject(subject string) (user models.User, err error) {
	if user, err = m.db.userOrmer.GetOneBySubject(subject); err != nil {
		return models.User{}, fmt.Errorf("cannot find user with subject %s", subject)
	}
	return
}

func (m *module) UserGetOneByUsername(username string) (user models.User, err error) {
	if user, err = m.db.userOrmer.GetOneByUsername(username); err != nil {
		return models.User{}, fmt.Errorf("cannot find user with username %s", username)
	}
	return
}
func (m *module) UserGetOneByEmail(email string) (user models.User, err error) {
	if user, err = m.db.userOrmer.GetOneByEmail(email); err != nil {
		return models.User{}, fmt.Errorf("cannot find user with email %s", email)
	}
	return
}

func (m *module) UserGetAll() (users []models.User, err error) {
	if users, err = m.db.userOrmer.GetAll(); err != nil {
		return []models.User{}, errors.New("cannot retrieve users")
	}
	return
}

func (m *module) UserUpdate(subject string, user datatransfers.UserUpdate) (err error) {
	if err = m.db.userOrmer.UpdateUser(models.User{
		Subject:    subject,
		GivenName:  user.GivenName,
		FamilyName: user.FamilyName,
		Email:      user.Email,
	}); err != nil {
		return errors.New("cannot update user")
	}
	return
}

func (m *module) UserUpdatePassword(subject, oldPassword, newPassword string) (err error) {
	var user models.User
	if user, err = m.db.userOrmer.GetOneBySubject(subject); err != nil {
		return errors2.ErrAuthIncorrectCredentials
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors2.ErrAuthIncorrectCredentials
	}
	var hashedPassword []byte
	if hashedPassword, err = bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost); err != nil {
		return errors.New("failed hashing password")
	}
	if err = m.db.userOrmer.UpdateUser(models.User{
		Subject:  subject,
		Password: string(hashedPassword),
	}); err != nil {
		return errors.New("cannot update user password")
	}
	return
}
