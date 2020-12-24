package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
)

func (m *module) GenerateAuthorizationCode(application models.Application) (authorizationCode string, err error) {
	authorizationCode = utils.GenerateRandomString(constants.AuthorizationCodeLength)
	if err = m.rd.SetEX(context.Background(), fmt.Sprintf("authorization_code::%s", application.ClientID),
		authorizationCode, constants.AuthorizationCodeExpiry).Err(); err != nil {
		return "", errors.New(fmt.Sprintf("failed storing authorization code. %v", err))
	}
	return
}
