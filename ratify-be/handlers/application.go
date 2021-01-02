package handlers

import (
	"errors"
	"fmt"
	"log"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
)

func (m *module) RegisterApplication(application datatransfers.ApplicationInfo, ownerSubject string) (clientID string, err error) {
	if application.Description == "" {
		application.Description = "New application"
	}
	if clientID, err = m.db.applicationOrmer.InsertApplication(models.Application{
		OwnerSubject: ownerSubject,
		ClientID:     utils.GenerateRandomString(constants.ClientIDLength),
		ClientSecret: utils.GenerateRandomString(constants.ClientSecretLength),
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

func (m *module) RetrieveAllApplications() (applications []models.Application, err error) {
	if applications, err = m.db.applicationOrmer.GetAll(); err != nil {
		return []models.Application{}, errors.New("cannot retrieve applications")
	}
	return
}

func (m *module) UpdateApplication(application datatransfers.ApplicationInfo) (err error) {
	if err = m.db.applicationOrmer.UpdateApplication(models.Application{
		ClientID:    application.ClientID,
		Name:        application.Name,
		Description: application.Description,
		LoginURL:    application.LoginURL,
		CallbackURL: application.CallbackURL,
		LogoutURL:   application.LogoutURL,
		Metadata:    application.Metadata,
	}); err != nil {
		return errors.New(fmt.Sprintf("error updating application. %v", err))
	}
	return
}

func (m *module) DeleteApplication(clientID string) (err error) {
	if err = m.db.applicationOrmer.DeleteApplication(clientID); err != nil {
		return errors.New(fmt.Sprintf("error deleting application. %v", err))
	}
	return
}
