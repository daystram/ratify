package handlers

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
)

func (m *module) ApplicationRegister(application datatransfers.ApplicationInfo, ownerSubject string) (clientID, clientSecret string, err error) {
	if application.Description == "" {
		application.Description = "New application"
	}
	clientSecret = utils.GenerateRandomString(constants.ClientSecretLength)
	var hashedClientSecret []byte
	if hashedClientSecret, err = bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost); err != nil {
		return "", "", errors.New("failed hashing client_secret")
	}
	if clientID, err = m.db.applicationOrmer.InsertApplication(models.Application{
		OwnerSubject: ownerSubject,
		ClientID:     utils.GenerateRandomString(constants.ClientIDLength),
		ClientSecret: string(hashedClientSecret),
		Name:         application.Name,
		Description:  application.Description,
		LoginURL:     application.LoginURL,
		CallbackURL:  application.CallbackURL,
		LogoutURL:    application.LogoutURL,
		Metadata:     application.Metadata,
	}); err != nil {
		return "", "", fmt.Errorf("error inserting application. %v", err)
	}
	return
}

func (m *module) ApplicationRenewClientSecret(clientID string) (clientSecret string, err error) {
	clientSecret = utils.GenerateRandomString(constants.ClientSecretLength)
	var hashedClientSecret []byte
	if hashedClientSecret, err = bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost); err != nil {
		return "", errors.New("failed hashing client_secret")
	}
	if err = m.db.applicationOrmer.UpdateApplication(models.Application{
		ClientID:     clientID,
		ClientSecret: string(hashedClientSecret),
	}); err != nil {
		return "", fmt.Errorf("error renewing application client_secret. %v", err)
	}
	return
}

func (m *module) ApplicationGetOneByClientID(clientID string) (application models.Application, err error) {
	if application, err = m.db.applicationOrmer.GetOneByClientID(clientID); err != nil {
		return models.Application{}, fmt.Errorf("cannot find application with client_id %s", clientID)
	}
	return
}

func (m *module) ApplicationGetOneByOwnerSubject(ownerSubject string) (applications []models.Application, err error) {
	if applications, err = m.db.applicationOrmer.GetAllByOwnerSubject(ownerSubject); err != nil {
		return []models.Application{}, errors.New("cannot retrieve applications")
	}
	return
}

func (m *module) ApplicationGetAll() (applications []models.Application, err error) {
	if applications, err = m.db.applicationOrmer.GetAll(); err != nil {
		return []models.Application{}, errors.New("cannot retrieve applications")
	}
	return
}

func (m *module) ApplicationUpdate(application datatransfers.ApplicationInfo) (err error) {
	if err = m.db.applicationOrmer.UpdateApplication(models.Application{
		ClientID:    application.ClientID,
		Name:        application.Name,
		Description: application.Description,
		LoginURL:    application.LoginURL,
		CallbackURL: application.CallbackURL,
		LogoutURL:   application.LogoutURL,
	}); err != nil {
		return fmt.Errorf("error updating application. %v", err)
	}
	return
}

func (m *module) ApplicationDelete(clientID string) (err error) {
	if err = m.db.applicationOrmer.DeleteApplication(clientID); err != nil {
		return fmt.Errorf("error deleting application. %v", err)
	}
	return
}
