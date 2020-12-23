package handlers

import (
	"errors"
	"fmt"
	"log"

	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/constants"
	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/datatransfers"
	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/models"
	"github.com/daystram/go-gin-gorm-boilerplate/ratify-be/utils"
)

func (m *module) RegisterApplication(application datatransfers.ApplicationInfo, ownerSubject string) (clientID string, err error) {
	if clientID, err = m.db.applicationOrmer.InsertApplication(models.Application{
		OwnerSubject: ownerSubject,
		ClientID:     utils.GenerateHexString(constants.ClientIDLength),
		ClientSecret: utils.GenerateHexString(constants.ClientSecretLength),
		Name:         application.Name,
		Description:  application.Description,
		LoginURL:     application.LoginURL,
		CallbackURL:  application.CallbackURL,
		LogoutURL:    application.LogoutURL,
		Metadata:     application.Metadata,
	}); err != nil {
		log.Print(err)
		return "", errors.New(fmt.Sprintf("error inserting application. %v", err))
	}
	return
}

func (m *module) RetrieveApplication(clientID string) (application models.Application, err error) {
	if application, err = m.db.applicationOrmer.GetOneByClientID(clientID); err != nil {
		return models.Application{}, errors.New(fmt.Sprintf("cannot find application with client_id %s", clientID))
	}
	return
}

// TODO: paginate
func (m *module) RetrieveOwnedApplications(ownerSubject string) (applications []models.Application, err error) {
	if applications, err = m.db.applicationOrmer.GetAllByOwnerSubject(ownerSubject); err != nil {
		return []models.Application{}, errors.New("cannot retrieve applications")
	}
	return
}

func (m *module) UpdateApplication(application datatransfers.ApplicationInfo) (err error) {
	if err = m.db.applicationOrmer.UpdateApplication(models.Application{
		Name:         application.Name,
		Description:  application.Description,
		LoginURL:     application.LoginURL,
		CallbackURL:  application.CallbackURL,
		LogoutURL:    application.LogoutURL,
		Metadata:     application.Metadata,
	}); err != nil {
		log.Print(err)
		return errors.New(fmt.Sprintf("error updating application. %v", err))
	}
	return
}
