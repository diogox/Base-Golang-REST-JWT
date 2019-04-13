package server

import "github.com/diogox/REST-JWT/server/pkg/models"

type EmailService interface {
	SendEmail(user models.User, opts models.NewEmail) error
}
