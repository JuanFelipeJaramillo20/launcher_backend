package mocks

import (
	"github.com/stretchr/testify/mock"
	"venecraft-back/cmd/entity"
)

type UserRoleRepository struct {
	mock.Mock
}

func (m *UserRoleRepository) AssignRole(userRole *entity.UserRole) error {
	args := m.Called(userRole)
	return args.Error(0)
}
