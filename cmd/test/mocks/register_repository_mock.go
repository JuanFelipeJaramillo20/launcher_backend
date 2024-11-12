package mocks

import (
	"github.com/stretchr/testify/mock"
	"venecraft-back/cmd/entity"
)

type RegisterRepository struct {
	mock.Mock
}

func (m *RegisterRepository) CreateRegister(register *entity.Register) error {
	args := m.Called(register)
	return args.Error(0)
}

func (m *RegisterRepository) GetAllRegisters() ([]entity.Register, error) {
	args := m.Called()
	return args.Get(0).([]entity.Register), args.Error(1)
}

func (m *RegisterRepository) GetRegisterByID(id uint64) (*entity.Register, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Register), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *RegisterRepository) DeleteRegister(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *RegisterRepository) UpdateRegister(register *entity.Register) error {
	args := m.Called(register)
	return args.Error(0)
}
