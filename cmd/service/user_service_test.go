package service

import (
	"errors"
	"testing"
	"time"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	roleRepo := new(mocks.RoleRepository)

	userService := NewUserService(userRepo, roleRepo)

	role := entity.Role{Name: "PLAYER"}
	user := &entity.User{
		FullName: "John Doe",
		Email:    "john@example.com",
		Nickname: "johndoe",
		Password: "Password123!",
	}

	// Mock expectations
	userRepo.On("GetUserByEmail", user.Email).Return(nil, nil)
	userRepo.On("GetUserByNickname", user.Nickname).Return(nil, nil)
	roleRepo.On("GetRoleByName", "PLAYER").Return(role, nil)
	userRepo.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(nil)

	err := userService.CreateUser(user, "PLAYER")

	assert.NoError(t, err)
	userRepo.AssertCalled(t, "CreateUser", mock.AnythingOfType("*entity.User"))
	roleRepo.AssertCalled(t, "GetRoleByName", "PLAYER")
}

func TestCreateUser_DuplicateEmailError(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	roleRepo := new(mocks.RoleRepository)

	userService := NewUserService(userRepo, roleRepo)

	existingUser := &entity.User{Email: "john@example.com"}
	user := &entity.User{
		FullName: "Jane Doe",
		Email:    "john@example.com",
		Nickname: "janedoe",
		Password: "password123",
	}

	userRepo.On("GetUserByEmail", user.Email).Return(existingUser, nil)

	err := userService.CreateUser(user, "PLAYER")

	assert.EqualError(t, err, "user with this email already exists")
	userRepo.AssertNotCalled(t, "CreateUser", user)
}

func TestRequestPasswordReset_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	emailClient := new(mocks.EmailClient)
	roleRepo := new(mocks.RoleRepository)

	userService := NewUserService(userRepo, roleRepo)

	user := &entity.User{Email: "john@example.com"}
	userRepo.On("GetUserByEmail", user.Email).Return(user, nil)
	userRepo.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(nil)
	emailClient.On("SendPasswordResetEmail", user.Email, mock.AnythingOfType("string")).Return(nil)

	err := userService.RequestPasswordReset(user.Email)

	assert.NoError(t, err)
	userRepo.AssertCalled(t, "UpdateUser", mock.AnythingOfType("*entity.User"))
	emailClient.AssertCalled(t, "SendPasswordResetEmail", user.Email, mock.AnythingOfType("string"))
}

func TestRequestPasswordReset_UserNotFound(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	emailClient := new(mocks.EmailClient)
	roleRepo := new(mocks.RoleRepository)

	userService := NewUserService(userRepo, roleRepo)

	userRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, errors.New("user not found"))

	err := userService.RequestPasswordReset("notfound@example.com")

	assert.EqualError(t, err, "user not found")
	userRepo.AssertNotCalled(t, "UpdateUser", mock.Anything)
	emailClient.AssertNotCalled(t, "SendPasswordResetEmail", mock.Anything, mock.Anything)
}

func TestResetPassword_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	roleRepo := new(mocks.RoleRepository)

	userService := NewUserService(userRepo, roleRepo)

	user := &entity.User{
		Email:                       "john@example.com",
		RecoverPasswordToken:        "valid-token",
		RecoverPasswordTokenExpires: time.Now().Add(15 * time.Minute),
	}
	userRepo.On("GetUserByResetToken", "valid-token").Return(user, nil)
	userRepo.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(nil)

	err := userService.ResetPassword("valid-token", "newpassword123")

	assert.NoError(t, err)
	userRepo.AssertCalled(t, "UpdateUser", mock.AnythingOfType("*entity.User"))
}

func TestResetPassword_InvalidToken(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	roleRepo := new(mocks.RoleRepository)

	userService := NewUserService(userRepo, roleRepo)

	userRepo.On("GetUserByResetToken", "invalid-token").Return(nil, errors.New("token not found"))

	err := userService.ResetPassword("invalid-token", "newpassword123")

	assert.EqualError(t, err, "invalid or expired token")
	userRepo.AssertNotCalled(t, "UpdateUser", mock.Anything)
}
