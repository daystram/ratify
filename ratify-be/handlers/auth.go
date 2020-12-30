package handlers

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
)

func (m *module) AuthenticateUser(credentials datatransfers.UserLogin) (user models.User, err error) {
	if user, err = m.db.userOrmer.GetOneByUsername(credentials.Username); err != nil {
		return models.User{}, errors.New("incorrect credentials")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return models.User{}, errors.New("incorrect credentials")
	}
	return user, nil
}

func (m *module) RegisterUser(userSignup datatransfers.UserSignup) (err error) {
	var hashedPassword []byte
	if hashedPassword, err = bcrypt.GenerateFromPassword([]byte(userSignup.Password), bcrypt.DefaultCost); err != nil {
		return errors.New("failed hashing password")
	}
	if _, err = m.db.userOrmer.InsertUser(models.User{
		GivenName:  userSignup.GivenName,
		FamilyName: userSignup.FamilyName,
		Username:   userSignup.Username,
		Email:      userSignup.Email,
		Password:   string(hashedPassword),
	}); err != nil  {
		return errors.New(fmt.Sprintf("error inserting user. %v", err))
	}
	return
}
