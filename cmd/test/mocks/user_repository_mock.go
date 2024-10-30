package mocks

import (
	"github.com/stretchr/testify/mock"
	"venecraft-back/cmd/entity"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) CreateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) GetAllUsers() ([]entity.User, error) {
	args := m.Called()
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *UserRepository) GetUserByID(id uint64) (*entity.User, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepository) UpdateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) DeleteUser(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepository) GetUserByNickname(nickname string) (*entity.User, error) {
	args := m.Called(nickname)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepository) GetUsersByRole(role string) ([]entity.User, error) {
	args := m.Called(role)
	return args.Get(0).([]entity.User), args.Error(1)
}

func (m *UserRepository) GetUserByResetToken(resetToken string) (*entity.User, error) {
	args := m.Called(resetToken)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepository) HasRole(id uint64, role string) bool {
	args := m.Called(id, role)
	return args.Bool(0)
}
