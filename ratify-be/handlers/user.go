package handlers

import (
	"errors"
	"fmt"

	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/datatransfers"
	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/models"
)

func (m *module) RetrieveUser(username string) (user models.User, err error) {
	if user, err = m.db.userOrmer.GetOneByUsername(username); err != nil {
		return models.User{}, errors.New(fmt.Sprintf("cannot find user with username %s", username))
	}
	return
}

func (m *module) UpdateUser(id string, user datatransfers.UserUpdate) (err error) {
	if err = m.db.userOrmer.UpdateUser(models.User{
		ID:    id,
		Email: user.Email,
	}); err != nil {
		return errors.New("cannot update user")
	}
	return
}
