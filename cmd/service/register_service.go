package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/enums"
	"venecraft-back/cmd/repository"
)

type RegisterService interface {
	CreateRegister(register *entity.Register) error
	ApproveRegister(id uint64) (*entity.User, error)
}

type registerService struct {
	registerRepo repository.RegisterRepository
	userRepo     repository.UserRepository
	roleRepo     repository.RoleRepository
	userRoleRepo repository.UserRoleRepository
}

func NewRegisterService(registerRepo repository.RegisterRepository, userRepo repository.UserRepository, roleRepo repository.RoleRepository, userRoleRepo repository.UserRoleRepository) RegisterService {
	return &registerService{registerRepo, userRepo, roleRepo, userRoleRepo}
}

func (s *registerService) CreateRegister(register *entity.Register) error {
	hashedPassword, err := hashPassword(register.Password)
	if err != nil {
		return err
	}
	register.Password = hashedPassword

	return s.registerRepo.CreateRegister(register)
}

func (s *registerService) ApproveRegister(id uint64) (*entity.User, error) {
	register, err := s.registerRepo.GetRegisterByID(id)
	if err != nil {
		return nil, errors.New("registration request not found")
	}

	user := &entity.User{
		FullName: register.FullName,
		Email:    register.Email,
		Nickname: register.Nickname,
		Password: register.Password,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	playerRole, err := s.roleRepo.GetRoleByName(enums.RolePlayer)
	if err != nil {
		return nil, errors.New("failed to assign role: PLAYER role not found")
	}

	userRole := &entity.UserRole{
		UserID: user.ID,
		RoleID: playerRole.ID,
	}
	err = s.userRoleRepo.AssignRole(userRole)
	if err != nil {
		return nil, err
	}

	err = s.registerRepo.DeleteRegister(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
