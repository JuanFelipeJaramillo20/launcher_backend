package mocks

import (
	"github.com/stretchr/testify/mock"
)

type EmailClient struct {
	mock.Mock
}

func (m *EmailClient) SendPasswordResetEmail(email string, link string) error {
	args := m.Called(email, link)
	return args.Error(0)
}
