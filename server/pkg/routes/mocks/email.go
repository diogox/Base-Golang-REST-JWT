package mocks

import "github.com/diogox/REST-JWT/server/pkg/models"

func NewMockEmailService() MockEmailService {
	return MockEmailService{}
}

type MockEmailService struct {}

func (MockEmailService) SendEmail(user models.User, opts models.NewEmail) error {
	return nil
}