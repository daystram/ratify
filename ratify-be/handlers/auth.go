package handlers

import (
	"errors"
	"fmt"
	"log"

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

func (m *module) RegisterUser(credentials datatransfers.UserSignup) (err error) {
	var hashedPassword []byte
	if hashedPassword, err = bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost); err != nil {
		return errors.New("failed hashing password")
	}
	if _, err = m.db.userOrmer.InsertUser(models.User{
		Username:  credentials.Username,
		Email:     credentials.Email,
		Password:  string(hashedPassword),
	}); err != nil {
		log.Print(err)
		return errors.New(fmt.Sprintf("error inserting user. %v", err))
	}
	return
}
