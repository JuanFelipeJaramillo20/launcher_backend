package service

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
	"venecraft-back/cmd/email"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/repository"
	"venecraft-back/cmd/utils"
)

type UserService interface {
	CreateUser(user *entity.User, roleName string) error
	GetAllUsers() ([]entity.User, error)
	GetUserByID(id uint64) (*entity.User, error)
	UpdateUser(id uint64, user *entity.User) error
	DeleteUser(id uint64) error
	RequestPasswordReset(email string) error
	ResetPassword(token string, newPassword string) error
}

type userService struct {
	userRepo    repository.UserRepository
	roleRepo    repository.RoleRepository
	emailClient *email.EmailClient
}

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) UserService {
	return &userService{userRepo, roleRepo, email.GetEmailClient()}
}

func (s *userService) CreateUser(user *entity.User, roleName string) error {
	if user.FullName == "" || user.Email == "" || user.Nickname == "" || user.Password == "" {
		return errors.New("all fields are required (FullName, Email, Nickname, Password)")
	}

	if !utils.IsValidEmail(user.Email) || !utils.IsValidNickname(user.Nickname) {
		return errors.New("invalid email or nickname format")
	}
	if err := utils.IsValidPassword(user.Password); err != nil {
		return err
	}

	if existingUser, _ := s.userRepo.GetUserByEmail(user.Email, false); existingUser != nil {
		return errors.New("user with this email already exists")
	}
	if existingUser, _ := s.userRepo.GetUserByNickname(user.Nickname); existingUser != nil {
		return errors.New("user with this nickname already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	if roleName != "" {
		role, err := s.roleRepo.GetRoleByName(roleName)
		if err != nil {
			return errors.New("specified role not found")
		}
		user.Roles = append(user.Roles, role)
	} else {
		role, err := s.roleRepo.GetRoleByName("PLAYER")
		if err != nil {
			return errors.New("default role PLAYER not found")
		}
		user.Roles = append(user.Roles, role)
	}

	// Save the user
	return s.userRepo.CreateUser(user)
}

func (s *userService) GetAllUsers() ([]entity.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *userService) GetUserByID(id uint64) (*entity.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) UpdateUser(userID uint64, userUpdate *entity.User) error {
	// Fetch the existing user from the database
	existingUser, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	if userUpdate.FullName != "" {
		existingUser.FullName = userUpdate.FullName
	}
	if userUpdate.Email != "" {
		existingUser.Email = userUpdate.Email
	}
	if userUpdate.Nickname != "" {
		existingUser.Nickname = userUpdate.Nickname
	}
	if userUpdate.Password != "" {
		hashedPassword, err := hashPassword(userUpdate.Password)
		if err != nil {
			return errors.New("failed to hash password")
		}
		existingUser.Password = hashedPassword
	}
	if userUpdate.RecoverPasswordToken != "" {
		existingUser.RecoverPasswordToken = userUpdate.RecoverPasswordToken
	}
	if !userUpdate.RecoverPasswordTokenExpires.IsZero() {
		existingUser.RecoverPasswordTokenExpires = userUpdate.RecoverPasswordTokenExpires
	}
	if userUpdate.Roles != nil {
		existingUser.Roles = userUpdate.Roles
	}
	existingUser.IsActive = userUpdate.IsActive

	return s.userRepo.UpdateUser(existingUser)
}

func (s *userService) DeleteUser(id uint64) error {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	user.IsActive = false
	return s.userRepo.UpdateUser(user)
}

func (s *userService) RequestPasswordReset(email string) error {
	user, err := s.userRepo.GetUserByEmail(email, false)
	if err != nil {
		return errors.New("user not found")
	}

	token, expiration, err := utils.GenerateResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %v", err)
	}

	user.RecoverPasswordToken = token
	user.RecoverPasswordTokenExpires = expiration
	if err := s.userRepo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to save reset token: %v", err)
	}

	return s.emailClient.SendPasswordResetEmail(user.Nickname, user.Email, token)
}

func (s *userService) ResetPassword(token, newPassword string) error {
	user, err := s.userRepo.GetUserByResetToken(token)
	if err != nil || time.Now().After(user.RecoverPasswordTokenExpires) {
		return errors.New("invalid or expired token")
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = hashedPassword
	user.RecoverPasswordToken = ""
	user.RecoverPasswordTokenExpires = time.Time{}

	return s.userRepo.UpdateUser(user)
}
