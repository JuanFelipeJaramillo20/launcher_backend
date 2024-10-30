package mocks

import (
	"github.com/stretchr/testify/mock"
	"venecraft-back/cmd/entity"
)

type RoleRepository struct {
	mock.Mock
}

func (m *RoleRepository) GetRoleByName(name string) (*entity.Role, error) {
	args := m.Called(name)
	return args.Get(0).(*entity.Role), args.Error(1)
}
