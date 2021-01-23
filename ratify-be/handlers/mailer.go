package handlers

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/matcornic/hermes/v2"

	"github.com/daystram/ratify/ratify-be/constants"
	"github.com/daystram/ratify/ratify-be/models"
	"github.com/daystram/ratify/ratify-be/utils"
)

func (m *module) MailerSendEmailVerification(user models.User) (err error) {
	token := utils.GenerateRandomString(constants.EmailVerificationTokenLength)
	if result := m.rd.SetEX(context.Background(), fmt.Sprintf(constants.RDTemVerificationToken, token),
		user.Subject, constants.EmailVerificationTokenExpiry); result.Err() != nil {
		return fmt.Errorf("invalid verification_token. %v", result.Err())
	}
	email := hermes.Email{
		Body: hermes.Body{
			Greeting:  "Hi",
			Signature: "Cheers!",
			Name:      user.Username,
			Intros: []string{
				"Welcome to Ratify!",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please click the following button to verify your email. This link expires in 30 minutes.",
					Button: hermes.Button{
						Color: "#00c3c3",
						Text:  "Confirm Email",
						Link:  fmt.Sprintf("https://ratify.daystram.com/verify?token=%s", token),
					},
				},
			},
		},
	}
	return m.sendEmail(user.Email, "Email Verification", email)
}

func (m *module) sendEmail(to, subject string, email hermes.Email) (err error) {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name:      "Ratify",
			Link:      "https://ratify.daystram.com/",
			Logo:      "https://ratify.daystram.com/images/logo.png",
			Copyright: fmt.Sprintf("Ratify %d. View on GitHub (https://github.com/daystram/ratify).", time.Now().Year()),
		},
	}

	var emailHTML, emailText string
	if emailHTML, err = h.GenerateHTML(email); err != nil {
		return fmt.Errorf("failed generating HTML email. %v", err)
	}
	if emailText, err = h.GeneratePlainText(email); err != nil {
		return fmt.Errorf("failed generating text email. %v", err)
	}

	from := mail.Address{
		Name:    constants.MailerIdentity,
		Address: "ratify@daystram.com",
	}
	message := gomail.NewMessage()
	message.SetHeader("From", from.String())
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", emailText)
	message.AddAlternative("text/html", emailHTML)
	return m.mailer.DialAndSend(message)
}
