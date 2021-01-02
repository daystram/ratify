package handlers

import (
	"errors"
	"fmt"

	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
)

func (m *module) RetrieveUserBySubject(subject string) (user models.User, err error) {
	if user, err = m.db.userOrmer.GetOneBySubject(subject); err != nil {
		return models.User{}, errors.New(fmt.Sprintf("cannot find user with subject %s", subject))
	}
	return
}

func (m *module) RetrieveUserByUsername(username string) (user models.User, err error) {
	if user, err = m.db.userOrmer.GetOneByUsername(username); err != nil {
		return models.User{}, errors.New(fmt.Sprintf("cannot find user with username %s", username))
	}
	return
}
func (m *module) RetrieveUserByEmail(email string) (user models.User, err error) {
	if user, err = m.db.userOrmer.GetOneByEmail(email); err != nil {
		return models.User{}, errors.New(fmt.Sprintf("cannot find user with email %s", email))
	}
	return
}

func (m *module) UpdateUser(id string, user datatransfers.UserUpdate) (err error) {
	if err = m.db.userOrmer.UpdateUser(models.User{
		Subject:    id,
		GivenName:  user.GivenName,
		FamilyName: user.FamilyName,
		Email:      user.Email,
	}); err != nil {
		return errors.New("cannot update user")
	}
	return
}
